<%!
import (
  "github.com/anden007/dp_clean_core/pkg"
)

type TableOption struct{
  RowNum int
  SearchSize int
  HideSearch bool
  ModalWidth string
  ApiName string
  Upload bool
  UploadThumb bool
  DaterangeSearch bool
  Password bool
  VueName string
  Api bool
  DefaultSort string
  DefaultSortType string
  LabelPosition string
  TotalRow int
  Span string
  FileUpload bool
  Editor bool
  Dict bool
  SearchDict bool
  CustomList bool
  SearchCustomList bool
}
%>
<%: func Table(fields []pkg.FormField, firstTwo []pkg.FormField, rest []pkg.FormField, options TableOption, buffer *bytes.Buffer) %>
<template>
  <div class="search">
    <Card>
      <%
      if options.SearchSize > 0 && !options.HideSearch{
      %>
      <Row v-show="openSearch" @keydown.enter.native="handleSearch">
        <Form ref="searchForm" :model="searchForm" inline :label-width="70">
        <%
        for _,item := range fields{
          if item.Searchable{
        %>
            <%
            if item.SearchType=="text"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <Input type="text" v-model="searchForm.<%==s item.Field%>.value" placeholder="请输入<%==s item.Name%>" clearable style="width: 200px">
                <Select
                  slot="prepend"
                  v-model="searchForm.<%==s item.Field%>.comparator"
                  transfer
                  style="width: 80px"
                >
                  <Option
                    v-for="(item, i) in getSearchComparator('string')"
                    :value="item.value"
                    :key="i"
                    >{{ item.key }}</Option
                  >
                </Select>
              </Input>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="select"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <Select v-model="searchForm.<%==s item.Field%>.value" placeholder="请选择" clearable style="width: 200px">
                <Option value="0">请自行编辑下拉菜单</Option>
              </Select>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="date"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <DatePicker type="date" v-model="searchForm.<%==s item.Field%>.value" placeholder="请选择" clearable style="width: 200px"></DatePicker>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="daterange"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <DatePicker :options="options" v-model="selectDate_<%==s item.Field%>" type="daterange" format="yyyy-MM-dd" clearable @on-change="selectDateRange_<%==s item.Field%>" placeholder="选择起始时间" style="width: 200px"></DatePicker>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="area"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <al-cascader v-model="searchForm.<%==s item.Field%>.value" data-type="name" level="<%==s item.SearchLevel%>" style="width:200px"/>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="dict"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <dict v-model="searchForm.<%==s item.Field%>.value" dict="<%==s item.SearchDictType%>" transfer style="width:200px"/>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="customList"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <customList v-model="searchForm.<%==s item.Field%>.value" url="<%==s item.SearchCustomUrl%>" transfer style="width:200px"/>
            </FormItem>
            <%
            }
            %>
        <%
          }
        }
        %>
          <FormItem style="margin-left:-35px;" class="br">
            <Button @click="handleSearch" type="primary" icon="ios-search">搜索</Button>
            <Button @click="handleReset">重置</Button>
          </FormItem>
        </Form>
      </Row>
      <%
      }
      %>
      <%
      if options.SearchSize>0 && options.HideSearch{
      %>
      <Row v-show="openSearch" @keydown.enter.native="handleSearch">
        <Form ref="searchForm" :model="searchForm" inline :label-width="70" class="search-form">
        <%
        for _, item := range firstTwo{
        %>
          <%
          if item.SearchType=="text"{
          %>
          <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
            <Input type="text" v-model="searchForm.<%==s item.Field%>.value" placeholder="请输入<%==s item.Name%>" clearable style="width: 200px">
                <Select
                  slot="prepend"
                  v-model="searchForm.<%==s item.Field%>.comparator"
                  transfer
                  style="width: 80px"
                >
                  <Option
                    v-for="(item, i) in getSearchComparator('string')"
                    :value="item.value"
                    :key="i"
                    >{{ item.key }}</Option
                  >
                </Select>
              </Input>
          </FormItem>
          <%
          }
          %>
          <%
          if item.SearchType=="select"{
          %>
          <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
            <Select v-model="searchForm.<%==s item.Field%>.value" placeholder="请选择" clearable style="width: 200px">
              <Option value="0">请自行编辑下拉菜单</Option>
            </Select>
          </FormItem>
          <%
          }
          %>
          <%
          if item.SearchType=="date"{
          %>
          <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
            <DatePicker type="date" v-model="searchForm.<%==s item.Field%>.value" placeholder="请选择" clearable style="width: 200px"></DatePicker>
          </FormItem>
          <%
          }
          %>
          <%
          if item.SearchType=="daterange"{
          %>
          <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
            <DatePicker :options="options" v-model="selectDate_<%==s item.Field%>" type="daterange" format="yyyy-MM-dd" clearable @on-change="selectDateRange_<%==s item.Field%>" placeholder="选择起始时间" style="width: 200px"></DatePicker>
          </FormItem>
          <%
          }
          %>
          <%
          if item.SearchType=="area"{
          %>
          <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
            <al-cascader v-model="searchForm.<%==s item.Field%>.value" data-type="name" level="<%==s item.SearchLevel%>" style="width:200px"/>
          </FormItem>
          <%
          }
          %>
          <%
          if item.SearchType=="dict"{
          %>
          <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
            <dict v-model="searchForm.<%==s item.Field%>.value" dict="<%==s item.SearchDictType%>" transfer style="width:200px"/>
          </FormItem>
          <%
          }
          %>
          <%
          if item.SearchType=="customList"{
          %>
          <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
            <customList v-model="searchForm.<%==s item.Field%>.value" url="<%==s item.SearchCustomUrl%>" transfer style="width:200px"/>
          </FormItem>
          <%
          }
          %>
        <%
        }
        %>
          <span v-if="drop">
          <%
          for _, item := range rest{
          %>
            <%
            if item.SearchType=="text"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <Input type="text" v-model="searchForm.<%==s item.Field%>.value" placeholder="请输入<%==s item.Name%>" clearable style="width: 200px"/>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="select"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <Select v-model="searchForm.<%==s item.Field%>.value" placeholder="请选择" clearable style="width: 200px">
                <Option value="0">请自行编辑下拉菜单</Option>
              </Select>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="date"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <DatePicker type="date" v-model="searchForm.<%==s item.Field%>.value" placeholder="请选择" clearable style="width: 200px"></DatePicker>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="daterange"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <DatePicker :options="options" v-model="selectDate_<%==s item.Field%>" type="daterange" format="yyyy-MM-dd" clearable @on-change="selectDateRange_<%==s item.Field%>" placeholder="选择起始时间" style="width: 200px"></DatePicker>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="area"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <al-cascader v-model="searchForm.<%==s item.Field%>.value" data-type="name" level="<%==s item.SearchLevel%>" style="width:200px"/>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="dict"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <dict v-model="searchForm.<%==s item.Field%>.value" dict="<%==s item.SearchDictType%>" transfer style="width:200px"/>
            </FormItem>
            <%
            }
            %>
            <%
            if item.SearchType=="customList"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <customList v-model="searchForm.<%==s item.Field%>.value" url="<%==s item.SearchCustomUrl%>" transfer style="width:200px"/>
            </FormItem>
            <%
            }
            %>
          <%
          }
          %>
          </span>
          <FormItem style="margin-left:-35px;" class="br">
            <Button @click="handleSearch" type="primary" icon="ios-search">搜索</Button>
            <Button @click="handleReset">重置</Button>
            <a class="drop-down" @click="dropDown">
              {{dropDownContent}}
              <Icon :type="dropDownIcon"></Icon>
            </a>
          </FormItem>
        </Form>
      </Row>
      <%
      }
      %>
      <Row class="operation">
        <ButtonGroup>
          <Button @click="add" type="primary" icon="md-add">添加</Button>
          <Button @click="delAll" icon="md-trash">批量删除</Button>
          <Button @click="getDataList" icon="md-refresh">刷新</Button>
          <% if options.SearchSize>0{ %>
          <Button icon="ios-search" @click="openSearch=!openSearch">{{openSearch ? "关闭搜索" : "开启搜索"}}</Button>
          <% } %>
        </ButtonGroup>
      </Row>
      <Alert show-icon v-show="openTip">
        已选择 <span class="select-count">{{selectList.length}}</span> 项
        <a class="select-clear" @click="clearSelectAll">清空</a>
      </Alert>
      <Table id="dataTable" :height="tableHeight" :loading="loading" border :columns="columns" :data="data" ref="table" sortable="custom" @on-sort-change="changeSort" @on-selection-change="changeSelect"></Table>
      <Row type="flex" justify="end" class="page">
        <Page :current="searchForm.pageNumber.value" :total="total" :page-size="searchForm.pageSize.value" @on-change="changePage" @on-page-size-change="changePageSize" :page-size-opts="[10,50,100]" size="small" show-total show-elevator show-sizer></Page>
      </Row>
    </Card>
    <Modal :title="modalTitle" v-model="modalVisible" :mask-closable='false' :width="<%==s options.ModalWidth%>">
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
              <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>" <% if item.Type=="upload"{ %>class="form-noheight"<% } %>>
                <%
                if item.Type=="text"{
                %>
                <Input v-model="form.<%==s item.Field%>" clearable/>
                <%
                }
                %>
                <%
                if item.Type=="textarea"{
                %>
                <Input v-model="form.<%==s item.Field%>" type="textarea" :rows="4" maxlength="250" show-word-limit/>
                <%
                }
                %>
                <%
                if item.Type=="select"{
                %>
                <Select v-model="form.<%==s item.Field%>" clearable>
                  <Option value="0">请自行编辑下拉菜单</Option>
                </Select>
                <%
                }
                %>
                <%
                if item.Type=="switch"{
                %>
                <i-switch v-model="form.<%==s item.Field%>"></i-switch>
                <%
                }
                %>
                <%
                if item.Type=="radio"{
                %>
                <RadioGroup v-model="form.<%==s item.Field%>">
                  <Radio label="0">请自行编辑单选框</Radio>
                  <Radio label="1">请自行编辑单选框</Radio>
                </RadioGroup>
                <%
                }
                %>
                <%
                if item.Type=="number"{
                %>
                <InputNumber v-model="form.<%==s item.Field%>" style="width:100%"></InputNumber>
                <%
                }
                %>
                <%
                if item.Type=="date"{
                %>
                <DatePicker type="date" v-model="form.<%==s item.Field%>" clearable style="width: 100%"></DatePicker>
                <%
                }
                %>
                <%
                if item.Type=="datetime"{
                %>
                <DatePicker type="datetime" v-model="form.<%==s item.Field%>" clearable transfer style="width: 100%"></DatePicker>
                <%
                }
                %>
                <%
                  if item.Type=="daterange"{
                %>
                <DatePicker type="daterange" v-model="form.<%==s item.Field%>" clearable style="width: 100%"></DatePicker>
                <%
                }
                %>
                <%
                if item.Type=="time"{
                %>
                <TimePicker type="time" v-model="form.<%==s item.Field%>" clearable style="width: 100%"></TimePicker>
                <%
                }
                %>
                <%
                if item.Type=="area"{
                %>
                <al-cascader v-model="form.<%==s item.Field%>" data-type="name" level="<%==s item.Level%>"/>
                <%
                }
                %>
                <%
                if item.Type=="slider"{
                %>
                <Slider v-model="form.<%==s item.Field%>"></Slider>
                <%
                }
                %>
                <%
                if item.Type=="upload"{
                %>
                <upload-pic-input v-model="form.<%==s item.Field%>"></upload-pic-input>
                <%
                }
                %>
                <%
                if item.Type=="uploadThumb"{
                %>
                <uploadThumb v-model="form.<%==s item.Field%>" multiple></uploadThumb>
                <%
                }
                %>
                <%
                if item.Type=="editor"{
                %>
                <editor id="editor-<%==i i%>-<%==i j%>" v-model="form.<%==s item.Field%>"></editor>
                <%
                }
                %>
                <%
                if item.Type=="password"{
                %>
                <password v-model="form.<%==s item.Field%>"></password>
                <%
                }
                %>
                <%
                if item.Type=="dict"{
                %>
                <dict v-model="form.<%==s item.Field%>" dict="<%==s item.DictType%>" transfer></dict>
                <%
                }
                %>
                <%
                if item.Type=="customList"{
                %>
                <customList v-model="form.<%==s item.Field%>" url="<%==s item.CustomUrl%>" transfer></customList>
                <%
                }
                %>
                <%
                if item.Type=="fileUpload"{
                %>
                <fileUpload v-model="form.<%==s item.Field%>"></fileUpload>
                <%
                }
                %>
              </FormItem>
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
      <div slot="footer">
        <Button type="text" @click="modalVisible=false">取消</Button>
        <Button type="primary" :loading="submitLoading" @click="handleSubmit">提交</Button>
      </div>
    </Modal>
  </div>
</template>

<script>
<%
if options.Api{
%>
// 根据你的实际请求api.js位置路径修改
import { get<%==s options.ApiName%>List, add<%==s options.ApiName%>, edit<%==s options.ApiName%>, delete<%==s options.ApiName%> } from "./api";
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
if options.Dict||options.SearchDict{
%>
import dict from "@/views/my-components/xboot/dict";
<%
}
%>
<%
if options.CustomList||options.SearchCustomList{
%>
import customList from "@/views/my-components/xboot/custom-list";
<%
}
%>
import { shortcuts } from "@/libs/shortcuts";
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
    if options.Dict||options.SearchDict{
    %>
    dict,
    <%
    }
    %>
    <%
    if options.CustomList||options.SearchCustomList{
    %>
    customList,
    <%
    }
    %>
  },
  data() {
    return {
      <% if options.SearchSize>0{ %>
      openSearch: true, // 显示搜索
      <% }%>
      openTip: true, // 显示提示
      loading: true, // 表单加载状态
      modalType: 0, // 添加或编辑标识
      modalVisible: false, // 添加或编辑显示
      modalTitle: "", // 添加或编辑标题
      <% if options.HideSearch { %>
      drop: false,
      dropDownContent: "展开",
      dropDownIcon: "ios-arrow-down",
      <% } %>
      initSearchForm: {},
      searchForm: { // 搜索框初始化对象
        pageNumber: { value: 1 }, // 当前页数
        pageSize: { value: 100 }, // 页面大小
        <%  if options.DefaultSort != "" {%>
        sort: { value: "<%==s options.DefaultSort%>" }, // 默认排序字段
        order: { value: "<%==s options.DefaultSortType%>" }, // 默认排序方式
        <%}%>
        <%
        for _, item := range fields{
        %>
        <%==s item.Field%>: { comparator: "==", value: "" },
        <%
        }
        %>
      },
      <% if options.DaterangeSearch{
        for _,item := range fields{
          if item.Searchable && item.SearchType=="daterange"{
      %>
      selectDate_<%==s item.Field%>: null,
      <% }} %>
      options: {
        shortcuts: shortcuts,
      },
      <% } %>
      form: { // 添加或编辑表单对象初始化数据
        <%
        for _,item := range fields{
          if item.Editable{
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
      // 表单验证规则
      formValidate: {
        <%
        for _,item := range fields{
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
      submitLoading: false, // 添加或编辑提交状态
      selectList: [], // 多选数据
      columns: [
        // 表头
        {
          type: "selection",
          width: 60,
          align: "center"
        },
        {
          type: "index",
          width: 60,
          align: "center"
        },
        <%
        for _, item := range fields{
          if item.TableShow{
        %>
        {
          title: "<%==s item.Name%>",
          key: "<%==s item.Field%>",
          width: 150,
          <%
          if item.Sortable{
          %>
          sortable: true,
          <%
          }else{
          %>
          sortable: false,
          <%
          }
          %>
          <%
          if item.DefaultSort{
          %>
          sortType: "<%==s item.DefaultSortType%>"
          <%
          }
          %>
        },
        <%
          }
        }
        %>
        {
          title: "操作",
          key: "action",
          align: "center",
          width: 150,
          fixed: "right",
          render: (h, params) => {
            return h("div", [
              h(
                "a",
                {
                  on: {
                    click: () => {
                      this.edit(params.row);
                    }
                  }
                },
                "编辑"
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
                    }
                  }
                },
                "删除"
              )
            ]);
          }
        }
      ],
      data: [], // 表单数据
      total: 0, // 表单数据总数
      tableHeight: 300,
    };
  },
  methods: {
    init() {
      this.initSearchForm = JSON.parse(JSON.stringify( this.searchForm ));
      this.getDataList();
      this.tableHeight = this.getTableHeight("#dataTable")
    },
    changePage(v) {
      this.searchForm.pageNumber.value = v;
      this.getDataList();
      this.clearSelectAll();
    },
    changePageSize(v) {
      this.searchForm.pageSize.value = v;
      this.getDataList();
    },
    <% if options.SearchSize>0{ %>
    handleSearch() {
      this.searchForm.pageNumber.value = 1;
      // this.searchForm.pageSize.value = 10;
      this.getDataList();
    },
    <% } %>
    handleReset() {
      this.searchForm = JSON.parse(JSON.stringify( this.initSearchForm ));
      <% if options.DaterangeSearch{
        for _,item := range fields{
          if item.Searchable && item.SearchType=="daterange"{
      %>
      this.selectDate_<%==s item.Field%> = null;
      <% }}} %>
      // 重新加载数据
      this.getDataList();
    },
    changeSort(e) {
      this.searchForm.sort.value = e.key;
      this.searchForm.order.value = e.order;
      if (e.order === "normal") {
        this.searchForm.order.value = "";
      }
      this.getDataList();
    },
    clearSelectAll() {
      this.$refs.table.selectAll(false);
    },
    changeSelect(e) {
      this.selectList = e;
    },
    <% if options.DaterangeSearch{
        for _,item := range fields{
          if item.Searchable && item.SearchType=="daterange"{
    %>
    selectDateRange_<%==s item.Field%>(v) {
      if (v) {
        this.searchForm.<%==s item.Field%>.comparator = "between";
        this.searchForm.<%==s item.Field%>.value = v[0] + " 00:00:00";
        this.searchForm.<%==s item.Field%>.value2 = v[1] + " 23:59:59";
      }
    },
    <% }}} %>
    <% if options.HideSearch{ %>
    dropDown() {
      if (this.drop) {
        this.dropDownContent = "展开";
        this.dropDownIcon = "ios-arrow-down";
      } else {
        this.dropDownContent = "收起";
        this.dropDownIcon = "ios-arrow-up";
      }
      this.drop = !this.drop;
    },
    <% } %>
    getDataList() {
      this.loading = true;
      <%
      if options.Api{
      %>
      get<%==s options.ApiName%>List(this.searchForm).then(res => {
        this.loading = false;
        if (res.success) {
          if(res.result.content){
            this.data = res.result.content;
            this.total = res.result.totalElements;
          }else{
            this.data = [];
            this.total = 0;
          }
          if (this.data.length == 0 && this.searchForm.pageNumber.value > 1) {
            this.searchForm.pageNumber.value = 0;
            this.getDataList();
          }
        }
      });
      <%
      } else {
      %>
      // 带多条件搜索参数获取表单数据 请自行修改接口
      // this.getRequest("请求路径", this.searchForm).then(res => {
      //   this.loading = false;
      //   if (res.success) {
      //     if(res.result.content){
      //       this.data = res.result.content;
      //       this.total = res.result.totalElements;
      //     }else{
      //       this.data = [];
      //       this.total = 0;
      //     }
      //     if (this.data.length == 0 && this.searchForm.pageNumber.value > 1) {
      //       this.searchForm.pageNumber.value = 0;
      //       this.getDataList();
      //     }
      //   }
      // });
      // 以下为模拟数据
      //this.data = [
      //];
      this.total = this.data.length;
      this.loading = false;
      <%
      }
      %>
    },
    handleSubmit() {
      this.$refs.form.validate(valid => {
        if (valid) {
          <%
          for _, item := range fields{
            if item.Editable&&item.Type=="date"{
          %>
          if (typeof this.form.<%==s item.Field%> == "object") {
            this.form.<%==s item.Field%> = this.format(this.form.<%==s item.Field%>, "yyyy-MM-dd HH:mm:ss");
          }
          <%
          }else if item.Editable&&item.Type=="datetime"{
          %>
          if (typeof this.form.<%==s item.Field%> == "object") {
            this.form.<%==s item.Field%> = this.format(this.form.<%==s item.Field%>, "yyyy-MM-dd HH:mm:ss");
          }
          <%
            }
          }
          %>
          this.submitLoading = true;
          if (this.modalType === 0) {
            // 添加 避免编辑后传入id等数据 记得删除
            delete this.form.id;
            if('createTime' in this.form){
              delete this.form.createTime;
            }
            <%
            if options.Api{
            %>
            add<%==s options.ApiName%>(this.form).then(res => {
              this.submitLoading = false;
              if (res.success) {
                this.$Message.success("操作成功");
                this.getDataList();
                this.modalVisible = false;
              }
            });
            <%
            } else {
            %>
            // this.postRequest("请求地址", this.form).then(res => {
            //   this.submitLoading = false;
            //   if (res.success) {
            //     this.$Message.success("操作成功");
            //     this.getDataList();
            //     this.modalVisible = false;
            //   }
            // });
            // 模拟请求成功
            this.submitLoading = false;
            this.$Message.success("操作成功");
            this.getDataList();
            this.modalVisible = false;
            <%
            }
            %>
          } else {
            // 编辑
            <%
            if options.Api{
            %>
            edit<%==s options.ApiName%>(this.form).then(res => {
              this.submitLoading = false;
              if (res.success) {
                this.$Message.success("操作成功");
                this.getDataList();
                this.modalVisible = false;
              }
            });
            <%
            } else {
            %>
            // this.postRequest("请求地址", this.form).then(res => {
            //   this.submitLoading = false;
            //   if (res.success) {
            //     this.$Message.success("操作成功");
            //     this.getDataList();
            //     this.modalVisible = false;
            //   }
            // });
            // 模拟请求成功
            this.submitLoading = false;
            this.$Message.success("操作成功");
            this.getDataList();
            this.modalVisible = false;
            <%
            }
            %>
          }
        }
      });
    },
    add() {
      this.modalType = 0;
      this.modalTitle = "添加";
      this.$refs.form.resetFields();
      delete this.form.id;
      this.modalVisible = true;
    },
    edit(v) {
      this.modalType = 1;
      this.modalTitle = "编辑";
      this.$refs.form.resetFields();
      // 转换null为""
      for (let attr in v) {
        if (v[attr] === null) {
          v[attr] = "";
        }
      }
      let str = JSON.stringify(v);
      let data = JSON.parse(str);
      this.form = data;
      this.modalVisible = true;
    },
    remove(v) {
      this.$Modal.confirm({
        title: "确认删除",
        // 记得确认修改此处
        content: "您确认要删除该条数据?",
        loading: true,
        onOk: () => {
          // 删除
          <%
          if options.Api{
          %>
          delete<%==s options.ApiName%>({ids: v.id}).then(res => {
            this.$Modal.remove();
            if (res.success) {
              this.$Message.success("操作成功");
              this.clearSelectAll();
              this.getDataList();
            }
          });
          <%
          } else {
          %>
          // this.deleteRequest("请求地址，如/deleteByIds/" + v.id).then(res => {
          //   this.$Modal.remove();
          //   if (res.success) {
          //     this.$Message.success("操作成功");
          //     this.clearSelectAll();
          //     this.getDataList();
          //   }
          // });
          // 模拟请求成功
          this.$Message.success("操作成功");
          this.clearSelectAll();
          this.$Modal.remove();
          this.getDataList();
          <%
          }
          %>
        }
      });
    },
    delAll() {
      if (this.selectList.length <= 0) {
        this.$Message.warning("您还未选择要删除的数据");
        return;
      }
      this.$Modal.confirm({
        title: "确认删除",
        content: "您确认要删除所选的 " + this.selectList.length + " 条数据?",
        loading: true,
        onOk: () => {
          let ids = "";
          this.selectList.forEach(function(e) {
            ids += e.id + ",";
          });
          ids = ids.substring(0, ids.length - 1);
          // 批量删除
          <%
          if options.Api{
          %>
          delete<%==s options.ApiName%>({ids: ids}).then(res => {
            this.$Modal.remove();
            if (res.success) {
              this.$Message.success("操作成功");
              this.clearSelectAll();
              this.getDataList();
            }
          });
          <%
          } else {
          %>
          // this.deleteRequest("请求地址，如/deleteByIds/" + ids).then(res => {
          //   this.$Modal.remove();
          //   if (res.success) {
          //     this.$Message.success("操作成功");
          //     this.clearSelectAll();
          //     this.getDataList();
          //   }
          // });
          // 模拟请求成功
          this.$Message.success("操作成功");
          this.$Modal.remove();
          this.clearSelectAll();
          this.getDataList();
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
// @import "@/styles/table-common.less";
.search {
    .operation {
        margin-bottom: 2vh;
    }
    .select-count {
        font-weight: 600;
        color: #40a9ff;
    }
    .select-clear {
        margin-left: 10px;
    }
    .page {
        margin-top: 2vh;
    }
    .drop-down {
        margin-left: 5px;
    }
}
</style>