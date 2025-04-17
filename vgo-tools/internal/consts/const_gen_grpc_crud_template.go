package consts

const TemplatePbEntityMessageContent = `

syntax = "proto3";

package {PackageName};

option go_package = "{GoPackage}";
{OptionContent}
{Imports}

{EntityMessage}

{RpcService}

{RpcServiceImplement}
`

const TemplateServiceContent = `
service {ServiceName}CrudService {
{Rpc}
}
`

const TemplateRpcServiceImplement = `
{RpcInvoke}

message {RpcRes} {
{ResField}
}
`
