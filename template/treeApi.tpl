<%: func TreeApi(apiName string, apiPath string, vueName string, buffer *bytes.Buffer) %>
// 统一请求路径前缀在libs/axios.js中修改
// import { getRequest, postBodyRequest, postRequest, deleteRequest } from '@/libs/axios';

// <%==s apiName%>获取一级数据
export const init<%==s apiName%> = (params) => {
    return getRequest('<%==s apiPath%>/getByParentId/0', params)
}
// <%==s apiName%>加载子级数据
export const load<%==s apiName%> = (id, params) => {
    return getRequest('<%==s apiPath%>/getByParentId/' + id, params)
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
    return postRequest('<%==s apiPath%>/delByIds/', params)
}
// <%==s apiName%>搜索
export const search<%==s apiName%> = (params) => {
    return getRequest('<%==s apiPath%>/search', params)
}