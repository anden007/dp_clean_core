<%: func Api(apiName string, apiPath string, vueName string, buffer *bytes.Buffer) %>
// 基础请求方法
import { getRequest, postRequest, putRequest, postBodyRequest, getNoAuthRequest, postNoAuthRequest, exportRequest } from "@/libs/axios";

// <%==s apiName%>分页获取数据
export const get<%==s apiName%>List = (params) => {
    return postBodyRequest('<%==s apiPath%>/getByCondition', params)
}
// <%==s apiName%>添加
export const add<%==s apiName%> = (params) => {
    return postBodyRequest('<%==s apiPath%>/add', params)
}
// <%==s apiName%>编辑
export const edit<%==s apiName%> = (params) => {
    return postBodyRequest('<%==s apiPath%>/edit', params)
}
// <%==s apiName%>删除
export const delete<%==s apiName%> = (params) => {
    return postRequest('<%==s apiPath%>/delByIds', params)
}