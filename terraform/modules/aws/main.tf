locals {
  # endpoints is a set of all of the distinct paths in var.endpoints
  endpoints = toset(var.functions.*.path)
  
  # methods is a map from method+path identifier strings to endpoint definitions
  methods = {
    for e in var.functions : 
    "${e.lambda_name}" => e
    if e.root_path == ""
  }
  proxy_methods = {
    for e in var.functions :
    "${e.lambda_name}" => e
    if e.root_path != ""
  }
}

data "archive_file" "lambda" {
  for_each = local.methods
  type = "zip"
  source_dir = format("./functions/%s", each.value.lambda_name)
  output_path = format("./functions/%s/dist/%s.zip", each.value.lambda_name, each.value.lambda_name)
}

resource "aws_api_gateway_rest_api" "api" {
  name = var.api_name
}

resource "aws_api_gateway_resource" "default" {
  for_each = local.endpoints

  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  path_part   = each.value
}

resource "aws_api_gateway_resource" "proxy" {
  for_each = local.endpoints

  rest_api_id = aws_api_gateway_rest_api.api.id
  parent_id   = aws_api_gateway_resource.default[each.value].id
  path_part   = "{proxy+}"
  depends_on = [
    aws_api_gateway_resource.default
  ]
}

resource "aws_api_gateway_method" "default" {
  for_each = local.methods

  rest_api_id = aws_api_gateway_resource.default[each.value.path].rest_api_id
  resource_id = aws_api_gateway_resource.default[each.value.path].id
  http_method = each.value.method
  authorization = length(each.value.authorization) > 0 ? each.value.authorization : "NONE"
}

resource "aws_api_gateway_method" "proxy" {
  for_each = local.proxy_methods

  rest_api_id = aws_api_gateway_resource.proxy[each.value.path].rest_api_id
  resource_id = aws_api_gateway_resource.proxy[each.value.path].id
  http_method = each.value.method
  authorization = length(each.value.authorization) > 0 ? each.value.authorization : "NONE"
  depends_on = [
    aws_api_gateway_method.default
  ]
}

resource "aws_api_gateway_integration" "default" {
  for_each = local.methods

  rest_api_id = aws_api_gateway_method.default[each.key].rest_api_id
  resource_id = aws_api_gateway_method.default[each.key].resource_id
  http_method = aws_api_gateway_method.default[each.key].http_method

  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = aws_lambda_function.default[each.key].invoke_arn
}

# Lambda
resource "aws_lambda_permission" "apigw_lambda" {
  for_each      = local.methods
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.default[each.key].function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"
}

resource "aws_lambda_function" "default" {
  for_each          = local.methods
  filename          = data.archive_file.lambda[each.value.lambda_name].output_path
  function_name     = each.value.lambda_name
  role              = aws_iam_role.role[each.key].arn
  handler           = each.value.lambda_name
  runtime           = each.value.lambda_runtime
  source_code_hash  = filebase64sha256(data.archive_file.lambda[each.value.lambda_name].output_path)
  # depends_on = [
  #   aws_iam_role_policy_attachment.lambda_logs,
  # ]
}

# IAM
resource "aws_iam_role" "role" {
  for_each = local.methods
  name = each.value.lambda_name

  assume_role_policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "lambda.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
      }
    ]
  }
  EOF
}

# resource "aws_cloudwatch_log_group" "default" {
#   for_each          = local.methods
#   name              = "/aws/lambda/${aws_lambda_function.default[each.value.lambda_name].function_name}"
#   retention_in_days = 7
# }

# See also the following AWS managed policy: AWSLambdaBasicExecutionRole
# resource "aws_iam_policy" "lambda_logging" {
#   name        = "lambda_logging"
#   path        = "/"
#   description = "IAM policy for logging from a lambda"

#   policy = <<EOF
# {
#   "Version": "2012-10-17",
#   "Statement": [
#     {
#       "Action": [
#         "logs:CreateLogGroup",
#         "logs:CreateLogStream",
#         "logs:PutLogEvents"
#       ],
#       "Resource": "arn:aws:logs:*:*:*",
#       "Effect": "Allow"
#     }
#   ]
# }
# EOF
# }

# resource "aws_iam_role_policy_attachment" "lambda_logs" {
#   for_each   = local.methods
#   role       = aws_iam_role.role[each.value.lambda_name].name
#   policy_arn = aws_iam_policy.lambda_logging[each.value.lambda_name].arn
# }

resource "aws_api_gateway_deployment" "default" {
  rest_api_id = aws_api_gateway_rest_api.api.id

  depends_on = [
    aws_api_gateway_resource.default,
    aws_api_gateway_method.default,
    aws_api_gateway_integration.default,
  ]
}


output "test" {
  value = local.proxy_methods
}