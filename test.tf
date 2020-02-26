# See https://github.com/hashicorp/hcl#syntax for the full reference
/* Different Types
of Comments */
key = "value"
valid_values_list = ["string", "number", "boolean", "object", "list"]
strings = "use_double_quotes"
really_long_strings = <<EOF
  can use heredocs
EOF
base_10_number = 1
octal_number = 02 # prefixed with 0
lists_with_objects = [{ key = "val"}, { key1 = false}]
boolean_list = [{ can_be = true}, { can_be = false}]
boolean_block_list {
  can_be = true
}
boolean_block_list {
  can_be = false
}
objects {
    can_be {
        nested = true
    }
}

# ex from terraform
resource "aws_route53_record" "site" {
    cname = "something.com"
}
variable "environment_name" {
    default =<<EOF
A long description about your environment
EOF
}

