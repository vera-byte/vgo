package consts

const TemplatePbEntityMessageContent = `
// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

syntax = "proto3";

package {PackageName};

option go_package = "{GoPackage}";
{OptionContent}
{Imports}

{EntityMessage}
`
