<%!
import (
  "github.com/anden007/dp_clean_core/pkg"
)

type TreeOption struct{
  RowNum int
  TotalRow int
  EnableTable bool
  TreeSpan string
  TreeEditSpan string
  LabelPosition string
  Span string
  IsDrawer bool
  ModalWidth string
  TotalRowTree int
  Api bool
  ApiName string
  FileUpload bool
  Upload bool
  UploadThumb bool
  Editor bool
  Password bool
  Dict bool
  CustomList bool
  VueName string
}
%>
<%: func Tree(fields []pkg.FormField, options TreeOption, buffer *bytes.Buffer) %>
<template>
  <div class="search">
    <Card>
      <Row class="operation" align="middle"<% if options.EnableTable{ %> justify="space-between"<% } %>>
        <% if options.EnableTable{ %><div><% } %>
        <Button @click="add" type="primary" icon="md-add" <% if options.EnableTable{ %>v-show="showType == 'tree'"<% } %>>添加子节点</Button>
        <Button @click="addRoot" icon="md-add">添加一级节点</Button>
        <Button @click="delAll" icon="md-trash">批量删除</Button>
        <Button @click="getParentList" icon="md-refresh">刷新</Button>
        <% if options.EnableTable{ %>
        <Input
          v-model="searchKey"
          suffix="ios-search"
          @on-change="search"
          placeholder="输入名称搜索"
          clearable
          style="width: 250px"
          v-show="showType == 'list'"
        />
        <% } %>
        <i-switch v-model="strict" size="large" style="margin-left:5px" <% if options.EnableTable{ %>v-show="showType == 'tree'"<% } %>>
          <span slot="open">级联</span>
          <span slot="close">单选</span>
        </i-switch>
        <% if options.EnableTable{ %>
        </div>
        <div>
          <RadioGroup v-model="showType" type="button">
            <Radio title="树结构" label="tree">
              <Icon type="md-list"></Icon>
            </Radio>
            <Radio title="列表" label="list">
              <Icon type="ios-apps"></Icon>
            </Radio>
          </RadioGroup>
        </div>
        <% } %>
      </Row>
      <Row type="flex" justify="start" :gutter="16" <% if options.EnableTable{ %>v-show="showType == 'tree'"<% } %>>
        <Col <%==s options.TreeSpan%>>
          <Alert show-icon>
            当前选择编辑：
            <span class="select-title">{{editTitle}}</span>
            <a class="select-clear" v-show="form.id && editTitle" @click="cancelEdit">取消选择</a>
          </Alert>
          <Input
            v-model="searchKey"
            suffix="ios-search"
            @on-change="search"
            placeholder="输入名称搜索"
            clearable
          />
          <div style="position: relative">
            <div class="tree-bar" :style="{maxHeight: maxHeight}">
              <Tree
                ref="tree"
                :data="data"
                :load-data="loadData"
                show-checkbox
                @on-check-change="changeSelect"
                @on-select-change="selectTree"
                :check-strictly="!strict"
              ></Tree>
            </div>
            <Spin size="large" fix v-if="loading"></Spin>
          </div>
        </Col>
        <Col <%==s options.TreeEditSpan%>>
          <Form ref="form" :model="form" :rules="formValidate" <% if options.LabelPosition=="left"{ %>:label-width="100"<% } %> label-position="<%==s options.LabelPosition%>">
            <%
            curr := 1
            for i:=0; i<options.TotalRow; i++ {
            %>
              <% if options.LabelPosition!="left"{ %>
              <Row :gutter="32">
              <% } %>
              <%
              for j:=0; j<options.RowNum; j++{
              if len(fields)==0||curr>len(fields){
                break
              }
              item := fields[curr-1];
              for{
                if !item.Editable && curr < len(fields){
                  curr++
                  item = fields[curr-1]
                }else{
                  break
                }
              }
              curr++;
              spanData := options.Span
              if item.Type=="editor"||item.Type=="textarea"{
                  spanData = "24";
              }
              %>
              <% if options.LabelPosition!="left"{ %>
              <Col span="<%==s spanData%>">
              <% } %>
              <%
              if item.Type=="xbootTreeChoose"{
              %>
            <FormItem label="上级节点" prop="parentTitle">
              <div style="display:flex;width:100%">
                <Input v-model="form.parentTitle" readonly style="margin-right:10px;"/>
                <Poptip transfer trigger="click" placement="right-start" title="选择上级节点" width="250">
                  <Button icon="md-list">选择节点</Button>
                  <div slot="content" class="tree-bar tree-select">
                    <Tree :data="dataEdit" :load-data="loadData" @on-select-change="selectTreeEdit"></Tree>
                    <Spin size="large" fix v-if="loadingEdit"></Spin>
                  </div>
                </Poptip>
              </div>
            </FormItem>
                <% } else if item.Type=="xbootTreeTitle"{ %>
            <FormItem label="名称" prop="title">
              <Input v-model="form.title"/>
            </FormItem>
                <% } else if item.Type=="xbootTreeSortOrder"{ %>
            <FormItem label="排序值" prop="sortOrder" class="block-tool">
              <Tooltip trigger="hover" placement="right" content="值越小越靠前，支持小数">
                <InputNumber :max="1000" :min="0" v-model="form.sortOrder" style="width:100%"></InputNumber>
              </Tooltip>
            </FormItem>
                <% }else{ %>
              <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>" <% if item.Type=="upload"{ %>class="form-noheight"<% } %>>
                <% if item.Type=="text"{ %>
                <Input v-model="form.<%==s item.Field%>" clearable/>
                <% } %>
                <% if item.Type=="textarea"{ %>
                <Input v-model="form.<%==s item.Field%>" type="textarea" :rows="4" maxlength="250" show-word-limit/>
                <% } %>
                <% if item.Type=="select"{ %>
                <Select v-model="form.<%==s item.Field%>" clearable transfer>
                  <Option value="0">请自行编辑下拉菜单</Option>
                </Select>
                <% } %>
                <% if item.Type=="switch"{ %>
                <i-switch v-model="form.<%==s item.Field%>"></i-switch>
                <% } %>
                <% if item.Type=="radio"{ %>
                <RadioGroup v-model="form.<%==s item.Field%>">
                  <Radio label="0">请自行编辑单选框</Radio>
                  <Radio label="1">请自行编辑单选框</Radio>
                </RadioGroup>
                <% } %>
                <% if item.Type=="number"{ %>
                <InputNumber v-model="form.<%==s item.Field%>" style="width:100%"></InputNumber>
                <% } %>
                <% if item.Type=="date"{ %>
                <DatePicker type="date" v-model="form.<%==s item.Field%>" clearable transfer style="width: 100%"></DatePicker>
                <% } %>
                <% if item.Type=="datetime"{ %>
                <DatePicker type="datetime" v-model="form.<%==s item.Field%>" clearable transfer style="width: 100%"></DatePicker>
                <% } %>
                <% if item.Type=="daterange"{ %>
                <DatePicker type="daterange" v-model="form.<%==s item.Field%>" clearable transfer style="width: 100%"></DatePicker>
                <% } %>
                <% if item.Type=="time"{ %>
                <TimePicker type="time" v-model="form.<%==s item.Field%>" clearable style="width: 100%" transfer></TimePicker>
                <% } %>
                <% if item.Type=="area"{ %>
                <al-cascader v-model="form.<%==s item.Field%>" data-type="name" level="<%==s item.Level%>" transfer/>
                <% } %>
                <% if item.Type=="slider"{ %>
                <Slider v-model="form.<%==s item.Field%>"></Slider>
                <% } %>
                <% if item.Type=="upload"{ %>
                <upload-pic-input v-model="form.<%==s item.Field%>"></upload-pic-input>
                <% } %>
                <% if item.Type=="uploadThumb"{ %>
                <uploadThumb v-model="form.<%==s item.Field%>" multiple></uploadThumb>
                <% } %>
                <% if item.Type=="editor"{ %>
                <editor id="editor-<%==i i%>-<%==i j%>" v-model="form.<%==s item.Field%>"></editor>
                <% } %>
                <% if item.Type=="password"{ %>
                <password v-model="form.<%==s item.Field%>"></password>
                <% } %>
                <% if item.Type=="dict"{ %>
                <dict v-model="form.<%==s item.Field%>" dict="<%==s item.DictType%>" transfer></dict>
                <% } %>
                <% if item.Type=="customList"{ %>
                <customList v-model="form.<%==s item.Field%>" url="<%==s item.CustomUrl%>" transfer></customList>
                <% } %>
                <% if item.Type=="fileUpload"{ %>
                <fileUpload v-model="form.<%==s item.Field%>"></fileUpload>
                <%
                }
                %>
              </FormItem>
              <%
              }
              %>
                <% if options.LabelPosition!="left"{ %>
                </Col>
                <% } %>
                <%
                }
                %>
              <% if options.LabelPosition!="left"{ %>
              </Row>
              <% } %>
              <%
              }
              %>
            <FormItem class="br">
              <Button
                @click="submitEdit"
                :loading="submitLoading"
                :disabled="!form.id || !editTitle"
                type="primary"
                icon="ios-create-outline"
              >修改并保存</Button>
              <Button @click="handleReset">重置</Button>
            </FormItem>
          </Form>
        </Col>
      </Row>
      <% if options.EnableTable{ %>
      <Alert show-icon v-show="showType == 'list'">
        已选择
        <span class="select-count">{{ selectList.length }}</span> 项
        <a class="select-clear" @click="clearSelectAll">清空</a>
      </Alert>
      <Table
        id="dataTable"
        :height="tableHeight"
        row-key="title"
        :load-data="loadData"
        :columns="columns"
        :data="data"
        :loading="loading"
        border
        :update-show-children="true"
        ref="table"
        @on-selection-change="showSelect"
        v-if="showType == 'list'"
      ></Table>
      <% } %>
    </Card>
    <% if !options.IsDrawer{ %>
    <Modal :title="modalTitle" v-model="modalVisible" :mask-closable="false" :width="<%==s options.ModalWidth%>">
    <% }else{ %>
    <Drawer :title="modalTitle" v-model="modalVisible" draggable :mask-closable="false":width="<%==s options.ModalWidth%>">
      <div :style="{ maxHeight: maxDrawerHeight }" class="drawer-content">
    <% } %>
      <Form ref="formAdd" :model="formAdd" :rules="formValidate" label-position="top">
        <div v-if="showParent">
          <FormItem label="上级节点：">{{form.title}}</FormItem>
        </div>
        <%
          var curr2 = 2;
          for i:=0; i<options.TotalRowTree; i++{
          %>
          <% if options.LabelPosition!="left"{ %>
          <Row :gutter="32">
          <% } %>
            <%
            for j:=0; j<options.RowNum; j++{
            if len(fields)==0||curr2>len(fields){
                break;
            }
            item := fields[curr2-1];
            for{
              if !item.Editable && curr2 <= len(fields){
                curr2++
                item = fields[curr2-1]
              }else{
                break
              }
            }
            curr2++;
            spanData := options.Span;
            if item.Type=="editor"||item.Type=="textarea"{
                spanData = "24";
            }
            %>
            <% if options.LabelPosition!="left"{ %>
            <Col span="<%==s spanData%>">
            <% } %>
            <%
            if item.Type=="xbootTreeTitle"{
            %>
             <FormItem label="名称" prop="title">
              <Input v-model="formAdd.title"/>
            </FormItem>
            <%
            } else if item.Type=="xbootTreeSortOrder"{
            %>
            <FormItem label="排序值" prop="sortOrder" class="block-tool">
              <Tooltip trigger="hover" placement="right" content="值越小越靠前，支持小数">
                <InputNumber :max="1000" :min="0" v-model="formAdd.sortOrder" style="width:100%"></InputNumber>
              </Tooltip>
            </FormItem>
            <%
            }else{
            %>
          <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>" <% if item.Type=="upload"{ %>class="form-noheight"<% } %>>
            <%
            if item.Type=="text"{
            %>
            <Input v-model="formAdd.<%==s item.Field%>" clearable/>
            <%
            }
            %>
            <%
            if item.Type=="textarea"{
            %>
            <Input v-model="formAdd.<%==s item.Field%>" type="textarea" :rows="4" maxlength="250" show-word-limit/>
            <%
            }
            %>
            <%
            if item.Type=="select"{
            %>
            <Select v-model="formAdd.<%==s item.Field%>" clearable transfer>
              <Option value="0">请自行编辑下拉菜单</Option>
            </Select>
            <%
            }
            %>
            <%
            if item.Type=="switch"{
            %>
            <i-switch v-model="formAdd.<%==s item.Field%>"></i-switch>
            <%
            }
            %>
            <%
            if item.Type=="radio"{
            %>
            <RadioGroup v-model="formAdd.<%==s item.Field%>">
              <Radio label="0">请自行编辑单选框</Radio>
              <Radio label="1">请自行编辑单选框</Radio>
            </RadioGroup>
            <%
            }
            %>
            <%
            if item.Type=="number"{
            %>
            <InputNumber v-model="formAdd.<%==s item.Field%>" style="width:100%"></InputNumber>
            <%
            }
            %>
            <%
            if item.Type=="date"{
            %>
            <DatePicker type="date" v-model="formAdd.<%==s item.Field%>" clearable transfer style="width: 100%"></DatePicker>
            <%
            }
            %>
            <%
            if item.Type=="datetime"{
            %>
            <DatePicker type="datetime" v-model="formAdd.<%==s item.Field%>" clearable transfer style="width: 100%"></DatePicker>
            <%
            }
            %>
            <%
              if item.Type=="daterange"{
            %>
            <DatePicker type="daterange" v-model="formAdd.<%==s item.Field%>" clearable transfer style="width: 100%"></DatePicker>
            <%
            }
            %>
            <%
            if item.Type=="time"{
            %>
            <TimePicker type="time" v-model="formAdd.<%==s item.Field%>" clearable transfer style="width: 100%"></TimePicker>
            <%
            }
            %>
            <%
            if item.Type=="area"{
            %>
            <al-cascader v-model="formAdd.<%==s item.Field%>" data-type="name" level="<%==s item.Level%>" transfer/>
            <%
            }
            %>
            <%
            if item.Type=="slider"{
            %>
            <Slider v-model="formAdd.<%==s item.Field%>"></Slider>
            <%
            }
            %>
            <%
            if item.Type=="upload"{
            %>
            <upload-pic-input v-model="formAdd.<%==s item.Field%>"></upload-pic-input>
            <%
            }
            %>
            <%
            if item.Type=="uploadThumb"{
            %>
            <uploadThumb v-model="formAdd.<%==s item.Field%>" multiple></uploadThumb>
            <%
            }
            %>
            <%
            if item.Type=="editor"{
            %>
            <editor id="editor-add-<%==i i%>-<%==i j%>" v-model="formAdd.<%==s item.Field%>"></editor>
            <%
            }
            %>
            <%
            if item.Type=="password"{
            %>
            <password v-model="formAdd.<%==s item.Field%>"></password>
            <%
            }
            %>
            <%
            if item.Type=="dict"{
            %>
            <dict v-model="formAdd.<%==s item.Field%>" dict="<%==s item.DictType%>" transfer></dict>
            <%
            }
            %>
            <%
            if item.Type=="customList"{
            %>
            <customList v-model="formAdd.<%==s item.Field%>" url="<%==s item.CustomUrl%>" transfer></customList>
            <%
            }
            %>
            <%
            if item.Type=="fileUpload"{
            %>
            <fileUpload v-model="formAdd.<%==s item.Field%>"></fileUpload>
            <%
            }
            %>
          </FormItem>
          <%
          }
          %>
            <% if options.LabelPosition!="left"{ %>
            </Col>
            <% } %>
            <%
            }
            %>
          <% if options.LabelPosition!="left"{ %>
          </Row>
          <% } %>
          <%
          }
          %>
      </Form>
      <% if options.IsDrawer{ %>
      </div>
      <% } %>
      <% if !options.IsDrawer{ %>
      <div slot="footer">
        <Button type="text" @click="modalVisible = false">取消</Button>
        <Button type="primary" :loading="submitLoading" @click="submitAdd">提交</Button>
      </div>
    </Modal>
      <% }else{ %>
      <div class="drawer-footer br">
        <Button type="primary" :loading="submitLoading" @click="submitAdd">提交</Button>
        <Button @click="modalVisible = false">取消</Button>
      </div>
    </Drawer>
      <% } %>
  </div>
</template>

<script>
<%
if options.Api{
%>
// 根据你的实际请求api.js位置路径修改
import { init<%==s options.ApiName%>, load<%==s options.ApiName%>, add<%==s options.ApiName%>, edit<%==s options.ApiName%>, delete<%==s options.ApiName%>, search<%==s options.ApiName%> } from "./api";
<%
}
%>
<%
if options.FileUpload{
%>
import fileUpload from "@/views/my-components/xboot/file-upload";
<%
}
%>
<%
if options.Upload{
%>
import uploadPicInput from "@/views/my-components/xboot/upload-pic-input";
<%
}
%>
<%
if options.UploadThumb{
%>
import uploadThumb from "@/views/my-components/xboot/upload-pic-thumb";
<%
}
%>
<%
if options.Editor{
%>
import editor from "@/views/my-components/xboot/editor";
<%
}
%>
<%
if options.Password{
%>
import password from "@/views/my-components/xboot/set-password";
<%
}
%>
<%
if options.Dict{
%>
import dict from "@/views/my-components/xboot/dict";
<%
}
%>
<%
if options.CustomList{
%>
import customList from "@/views/my-components/xboot/custom-list";
<%
}
%>
export default {
  name: "<%==s options.VueName%>",
  components: {
    <%
    if options.FileUpload{
    %>
      fileUpload,
    <%
    }
    %>
    <%
    if options.Upload{
    %>
    uploadPicInput,
    <%
    }
    %>
    <%
    if options.UploadThumb{
    %>
    uploadThumb,
    <%
    }
    %>
    <%
    if options.Editor{
    %>
    editor,
    <%
    }
    %>
    <%
    if options.Password{
    %>
    password,
    <%
    }
    %>
    <%
    if options.Dict{
    %>
    dict,
    <%
    }
    %>
    <%
    if options.CustomList{
    %>
    customList,
    <%
    }
    %>
  },
  data() {
    return {
      <% if options.EnableTable{ %>
      showType: "tree",
      <% } %>
      <% if options.IsDrawer{ %>
      maxDrawerHeight: 510,
      <% } %>
      maxHeight: "500px",
      strict: true,
      loading: false, // 树加载状态
      loadingEdit: false, // 编辑上级树加载状态
      modalVisible: false, // 添加显示
      selectList: [], // 多选数据
      showParent: false, // 显示上级标识
      modalTitle: "", // 添加标题
      editTitle: "", // 编辑节点名称
      searchKey: "", // 搜索树
      form: {
        // 编辑对象初始化数据
        id: "",
        title: "",
        parentId: "",
        parentTitle: "",
        sortOrder: 0,
        <%
        for _, item := range fields{
          if item.Editable {
        %>
        <% if item.Type=="switch"{ %>
        <%==s item.Field%>: true,
        <% }else if item.Type=="number"||item.Type=="slider"{ %>
        <%==s item.Field%>: 0,
        <% }else if item.Type=="daterange"||item.Type=="area"||item.Type=="uploadThumb"{ %>
        <%==s item.Field%>: [],
        <% }else{ %>
        <%==s item.Field%>: "",
        <% } %>
        <%
          }
        }
        %>
      },
      formAdd: {
        // 添加对象初始化数据
      },
      formValidate: {
        // 表单验证规则
        title: [{ required: true, message: "不能为空", trigger: "change" }],
        sortOrder: [
          {
            required: true,
            type: "number",
            message: "排序值不能为空",
            trigger: "change"
          }
        ],
        <%
        for _, item := range fields{
          if item.Editable&&item.Validate{
        %>
        <% if item.Type=="area"||item.Type=="uploadThumb"{ %>
        <%==s item.Field%>: [{ type: "array", required: true, message: "不能为空", trigger: "change" }],
        <% }else if item.Type=="date"||item.Type=="datetime"{ %>
        <%==s item.Field%>: [{ type: "date", required: true, message: "不能为空", trigger: "change" }],
        <% }else if item.Type=="daterange"{ %>
        <%==s item.Field%>: [{ type: "array", required: true, trigger: "change", fields: { 0: { type: "date", message: "起始日期不能为空", required: true }, 1: { type: "date", message: "结束日期不能为空", required: true } }, }],
        <% }else if item.Type=="number"||item.Type=="slider"{ %>
        <%==s item.Field%>: [{ type: "number", required: true, message: "不能为空", trigger: "change" }],
        <% }else if item.Type=="switch"{ %>
        <%==s item.Field%>: [{ type: "boolean", required: true, message: "不能为空", trigger: "change" }],
        <% }else{ %>
        <%==s item.Field%>: [{ required: true, message: "不能为空", trigger: "change" }],
        <% } %>
        <%
          }
        }
        %>
      },
      submitLoading: false,
      data: [],
      dataEdit: [],
      <% if options.EnableTable{ %>
      columns: [
        {
          type: "selection",
          width: 60,
          align: "center",
        },
        {
          type: "index",
          width: 60,
          align: "center",
        },
        {
          title: "名称",
          key: "title",
          minWidth: 120,
          sortable: true,
          tree: true,
        },
        {
          title: "排序",
          key: "sortOrder",
          width: 150,
          sortable: true,
          align: "center",
          sortType: "asc",
        },
        {
          title: "创建时间",
          key: "createTime",
          sortable: true,
          width: 200,
        },
        {
          title: "操作",
          key: "action",
          width: 300,
          align: "center",
          render: (h, params) => {
            return h("div", [
              h(
                "a",
                {
                  on: {
                    click: () => {
                      this.tableAdd(params.row);
                    },
                  },
                },
                "添加子节点"
              ),
              h("Divider", {
                props: {
                  type: "vertical",
                },
              }),
              h(
                "a",
                {
                  on: {
                    click: () => {
                      this.remove(params.row);
                    },
                  },
                },
                "删除"
              ),
            ]);
          },
        },
      ],
      <% } %>
      tableHeight: 300,
    };
  },
  methods: {
    init() {
      // 计算高度
      let height = document.documentElement.clientHeight;
      this.maxHeight = Number(height - 287) + "px";
      <% if options.IsDrawer{ %>
      this.maxDrawerHeight = Number(height - 121) + "px";
      <% } %>
      // 初始化一级节点
      this.getParentList();
      // 初始化一级节点为编辑上级节点使用
      this.getParentListEdit();
      this.tableHeight = this.getTableHeight("#dataTable")
    },
    getParentList() {
      <%
      if options.Api{
      %>
      this.loading = true;
      init<%==s options.ApiName%>().then(res => {
        this.loading = false;
        if (res.success) {
          res.result.forEach(function(e) {
            if (e.isParent) {
              e.loading = false;
              e.children = [];
              <% if options.EnableTable{ %>
              e._loading = false;
              <% } %>
            }
          });
          this.data = res.result;
        }
      });
      <%
      } else {
      %>
      // this.loading = true;
      // this.getRequest("一级数据请求路径，如/tree/getByParentId/0").then(res => {
      //   this.loading = false;
      //   if (res.success) {
      //     res.result.forEach(function(e) {
      //       if (e.isParent) {
      //         e.loading = false;
      //         e.children = [];
      <% if options.EnableTable{ %>
      //         e._loading = false;
      <% } %>
      //       }
      //     });
      //     this.data = res.result;
      //   }
      // });
      // 模拟请求成功
      this.data = [
      ];
      <%
      }
      %>
    },
    getParentListEdit() {
      <%
      if options.Api{
      %>
      this.loadingEdit = true;
      init<%==s options.ApiName%>().then(res => {
        this.loadingEdit = false;
        if (res.success) {
          res.result.forEach(function(e) {
            if (e.isParent) {
              e.loading = false;
              e.children = [];
            }
          });
          // 头部加入一级
          let first = {
            id: "0",
            title: "一级节点"
          };
          res.result.unshift(first);
          this.dataEdit = res.result;
        }
      });
      <%
      } else {
      %>
      // this.loadingEdit = true;
      // this.getRequest("/tree/getByParentId/0").then(res => {
      //   this.loadingEdit = false;
      //   if (res.success) {
      //     res.result.forEach(function(e) {
      //       if (e.isParent) {
      //         e.loading = false;
      //         e.children = [];
      //       }
      //     });
      //     // 头部加入一级
      //     let first = {
      //       id: "0",
      //       title: "一级节点"
      //     };
      //     res.result.unshift(first);
      //     this.dataEdit = res.result;
      //   }
      // });
      // 模拟请求成功
      this.dataEdit = [
      ];
      <%
      }
      %>
    },
    loadData(item, callback) {
      <%
      if options.Api{
      %>
      load<%==s options.ApiName%>(item.id).then(res => {
        if (res.success) {
          res.result.forEach(function(e) {
            if (e.isParent) {
              e.loading = false;
              e.children = [];
              <% if options.EnableTable{ %>
              e._loading = false;
              <% } %>
            }
          });
          callback(res.result);
        }
      });
      <%
      } else {
      %>
      // 异步加载树子节点数据
      // this.getRequest("请求路径，如/tree/getByParentId/" + item.id).then(res => {
      //   if (res.success) {
      //     res.result.forEach(function(e) {
      //       if (e.isParent) {
      //         e.loading = false;
      //         e.children = [];
      <% if options.EnableTable{ %>
      //         e._loading = false;
      <% } %>
      //       }
      //     });
      //     callback(res.result);
      //   }
      // });
      <%
      }
      %>
    },
    search() {
      // 搜索树
      if (this.searchKey) {
        <%
        if options.Api{
        %>
        this.loading = true;
        search<%==s options.ApiName%>({ title: this.searchKey }).then(res => {
          this.loading = false;
          if (res.success) {
            res.result.forEach(function(e) {
              if (e.isParent) {
                e.loading = false;
                e.children = [];
                <% if options.EnableTable{ %>
                e._loading = false;
                <% } %>
              }
            });
            this.data = res.result;
          }
        });
        <%
        } else {
        %>
        // 模拟请求
        // this.loading = true;
        // this.getRequest("搜索请求路径", { title: this.searchKey }).then(res => {
        //   this.loading = false;
        //   if (res.success) {
        //     res.result.forEach(function(e) {
        //       if (e.isParent) {
        //         e.loading = false;
        //         e.children = [];
        <% if options.EnableTable{ %>
        //         e._loading = false;
        <% } %>
        //       }
        //     });
        //     this.data = res.result;
        //   }
        // });
        // 模拟请求成功
        this.data = [
        ];
        <%
        }
        %>
      } else {
        // 为空重新加载
        this.getParentList();
      }
    },
    selectTree(v) {
      if (v.length > 0) {
        // 转换null为""
        for (let attr in v[0]) {
          if (v[0][attr] === null) {
            v[0][attr] = "";
          }
        }
        let str = JSON.stringify(v[0]);
        let data = JSON.parse(str);
        this.form = data;
        this.editTitle = data.title;
      } else {
        this.cancelEdit();
      }
    },
    cancelEdit() {
      let data = this.$refs.tree.getSelectedNodes()[0];
      if (data) {
        data.selected = false;
      }
      this.$refs.form.resetFields();
      this.form.id = "";
      this.editTitle = "";
    },
    selectTreeEdit(v) {
      if (v.length > 0) {
        // 转换null为""
        for (let attr in v[0]) {
          if (v[0][attr] === null) {
            v[0][attr] = "";
          }
        }
        let str = JSON.stringify(v[0]);
        let data = JSON.parse(str);
        if (this.form.id == data.id) {
          this.$Message.warning("请勿选择自己作为父节点");
          v[0].selected = false;
          return;
        }
        this.form.parentId = data.id;
        this.form.parentTitle = data.title;
      }
    },
    handleReset() {
      this.$refs.form.resetFields();
      this.form.status = 0;
    },
    submitEdit() {
      this.$refs.form.validate(valid => {
        if (valid) {
          if (!this.form.id) {
            this.$Message.warning("请先点击选择要修改的节点");
            return;
          }
          <%
          for _, item := range fields{
            if(item.Editable&&item.Type=="date"){
          %>
          if (typeof this.form.<%==s item.Field%> == "object") {
            this.form.<%==s item.Field%> = this.format(this.form.<%==s item.Field%>, "yyyy-MM-dd HH:mm:ss");
          }
          <%
          }else if(item.Editable&&item.Type=="datetime"){
          %>
          if (typeof this.form.<%==s item.Field%> == "object") {
            this.form.<%==s item.Field%> = this.format(this.form.<%==s item.Field%>, "yyyy-MM-dd HH:mm:ss");
          }
          <%
            }
          }
          %>
          this.submitLoading = true;
          <%
          if options.Api{
          %>
          edit<%==s options.ApiName%>(this.form).then(res => {
          this.submitLoading = false;
          if (res.success) {
            this.$Message.success("编辑成功");
              this.init();
              this.modalVisible = false;
            }
          });
          <%
          } else {
          %>
          // this.postRequest("请求路径，如/tree/edit", this.form).then(res => {
          //   this.submitLoading = false;
          //   if (res.success) {
          //     this.$Message.success("编辑成功");
          //     this.init();
          //     this.modalVisible = false;
          //   }
          // });
          // 模拟成功
          this.submitLoading = false;
          this.$Message.success("编辑成功");
          this.modalVisible = false;
          <%
          }
          %>
        }
      });
    },
    submitAdd() {
      this.$refs.formAdd.validate(valid => {
        if (valid) {
          this.submitLoading = true;
          <%
          if options.Api{
          %>
          add<%==s options.ApiName%>(this.formAdd).then(res => {
            this.submitLoading = false;
            if (res.success) {
              this.$Message.success("添加成功");
              this.init();
              this.modalVisible = false;
            }
          });
          <%
          } else {
          %>
          // this.postRequest("请求路径，如/tree/add", this.formAdd).then(res => {
          //   this.submitLoading = false;
          //   if (res.success) {
          //     this.$Message.success("添加成功");
          //     this.init();
          //     this.modalVisible = false;
          //   }
          // });
          // 模拟成功
          this.submitLoading = false;
          this.$Message.success("添加成功");
          this.modalVisible = false;
          <%
          }
          %>
        }
      });
    },
    add() {
      if (this.form.id == "" || this.form.id == null) {
        this.$Message.warning("请先点击选择一个节点");
        return;
      }
      this.modalTitle = "添加子节点";
      this.showParent = true;
      if (!this.form.children) {
        this.form.children = [];
      }
      this.formAdd = {
        parentId: this.form.id,
        sortOrder: this.form.children.length + 1
      };
      this.modalVisible = true;
    },
    addRoot() {
      this.modalTitle = "添加一级节点";
      this.showParent = false;
      this.formAdd = {
        parentId: "",
        sortOrder: this.data.length + 1
      };
      this.modalVisible = true;
    },
    changeSelect(v) {
      this.selectList = v;
    },
    <% if options.EnableTable{ %>
    clearSelectAll() {
      this.$refs.table.selectAll(false);
    },
    tableAdd(v) {
      this.form = v;
      this.add();
      this.editTitle = "";
      let data = this.$refs.tree.getSelectedNodes()[0];
      if (data) {
        data.selected = false;
      }
    },
    showSelect(e) {
      this.selectList = e;
    },
    remove(v) {
      this.selectList = [];
      this.selectList.push(v);
      this.delAll();
    },
    <% } %>
    delAll() {
      if (this.selectList.length <= 0) {
        this.$Message.warning("您还未勾选要删除的数据");
        return;
      }
      this.$Modal.confirm({
        title: "确认删除",
        content: "您确认要删除所选的 " + this.selectList.length + " 条数据及其下级所有数据?",
        loading: true,
        onOk: () => {
          let ids = "";
          this.selectList.forEach(function(e) {
            ids += e.id + ",";
          });
          ids = ids.substring(0, ids.length - 1);
           <%
          if options.Api{
          %>
          delete<%==s options.ApiName%>({ids: ids}).then(res => {
            this.$Modal.remove();
            if (res.success) {
              this.$Message.success("删除成功");
              this.selectList = [];
              this.cancelEdit();
              this.init();
            }
          });
          <%
          } else {
          %>
          // this.deleteRequest("请求路径，如/tree/delByIds/" + ids).then(res => {
          //   this.$Modal.remove();
          //   if (res.success) {
          //     this.$Message.success("删除成功");
          //     this.selectList = [];
          //     this.cancelEdit();
          //     this.init();
          //   }
          // });
          // 模拟成功
          this.$Modal.remove();
          this.$Message.success("删除成功");
          this.selectList = [];
          this.cancelEdit();
          <%
          }
          %>
        }
      });
    }
  },
  mounted() {
    this.init();
  }
};
</script>
<style lang="less">
// 建议引入通用样式 具体路径自行修改 可删除下面样式代码
// @import "@/styles/tree-common.less";
// @import "@/styles/drawer-common.less";
.search {
    .operation {
        margin-bottom: 2vh;
    }
    .select-title {
        font-weight: 600;
        color: #40a9ff;
    }
    .select-clear {
        margin-left: 10px;
    }
}

.tree-bar {
    overflow: auto;
    margin-top: 5px;
    position: relative;
    min-height: 80px;
}

.tree-select {
    max-height: 500px;
}

.tree-bar::-webkit-scrollbar {
    width: 6px;
    height: 6px;
}

.tree-bar::-webkit-scrollbar-thumb {
    border-radius: 4px;
    -webkit-box-shadow: inset 0 0 2px #d1d1d1;
    background: #e4e4e4;
}

.drawer-footer {
    z-index: 10;
    width: 100%;
    position: absolute;
    bottom: 0;
    left: 0;
    border-top: 1px solid #e8e8e8;
    padding: 10px 16px;
    text-align: left;
    background: #fff;
}

.drawer-content {
    overflow: auto;
}

.drawer-content::-webkit-scrollbar {
    display: none;
}

.drawer-header {
    display: flex;
    align-items: center;
    margin-bottom: 16px;
    font-size: 16px;
    color: rgba(0, 0, 0, 0.85);
}

.drawer-title {
    font-size: 16px;
    color: rgba(0, 0, 0, 0.85);
    display: block;
    margin-bottom: 16px;
}
</style>