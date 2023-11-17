<%!
import (
    "github.com/anden007/dp_clean_core/pkg"
)

type TableIndexOption struct{
    SearchSize int
    HideSearch bool
    ApiName string
    DaterangeSearch bool
    Api bool
    DefaultSort string
    DefaultSortType string    
    SearchDict bool    
    SearchCustomList bool
    VueName string
}
%>
<%: func TableIndex(fields []pkg.FormField, firstTwo []pkg.FormField, rest []pkg.FormField, options TableIndexOption, buffer *bytes.Buffer) %>
<template>
  <div class="search">
    <add v-if="currView=='add'" @close="currView='index'" @submited="submited" />
    <edit v-if="currView=='edit'" @close="currView='index'" @submited="submited" :data="formData" />
    <Card v-show="currView=='index'">
        <%
        if options.SearchSize>0 && !options.HideSearch{
        %>
        <Row <% if options.SearchSize>0{ %>v-show="openSearch"<% } %> @keydown.enter.native="handleSearch">
          <Form ref="searchForm" :model="searchForm" inline :label-width="70">
          <%
          for _, item := range fields{
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
              if item.SearchType=="date" {
              %>
              <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
                <DatePicker type="date" v-model="searchForm.<%==s item.Field%>.value" placeholder="请选择" clearable style="width: 200px"></DatePicker>
              </FormItem>
              <%
              }
              %>
              <%
              if item.SearchType=="daterange" {
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
        <Row @keydown.enter.native="handleSearch">
          <Form ref="searchForm" :model="searchForm" inline :label-width="70" class="search-form">
          <%
          for _, item := range firstTwo{
          %>
            <%
            if item.SearchType=="text"{
            %>
            <FormItem label="<%==s item.Name%>" prop="<%==s item.Field%>">
              <Input type="text" v-model="searchForm.<%==s item.Field%>.value" placeholder="请输入<%==s item.Name%>" clearable style="width: 200px"
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
                <Input type="text" v-model="searchForm.<%==s item.Field%>.value" placeholder="请输入<%==s item.Name%>" clearable style="width: 200px"
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
        <Button @click="add" type="primary" icon="md-add">添加</Button>
        <Button @click="delAll" icon="md-trash">批量删除</Button>
        <Button @click="getDataList" icon="md-refresh">刷新</Button>
        <Button type="dashed" @click="openTip=!openTip">{{openTip ? "关闭提示" : "开启提示"}}</Button>
      </Row>
      <Alert show-icon v-show="openTip">
        已选择
        <span class="select-count">{{selectList.length}}</span> 项
        <a class="select-clear" @click="clearSelectAll">清空</a>
      </Alert>
      <Table
        id="dataTable"
        :height="tableHeight"
        :loading="loading"
        border
        :columns="columns"
        :data="data"
        ref="table"
        sortable="custom"
        @on-sort-change="changeSort"
        @on-selection-change="changeSelect"
      ></Table>
      <Row type="flex" justify="end" class="page">
        <Page
          :current="searchForm.pageNumber.value"
          :total="total"
          :page-size="searchForm.pageSize.value"
          @on-change="changePage"
          @on-page-size-change="changePageSize"
          :page-size-opts="[10,20,50]"
          size="small"
          show-total
          show-elevator
          show-sizer
        ></Page>
      </Row>
    </Card>
  </div>
</template>

<script>
<%
if options.Api{
%>
// 根据你的实际请求api.js位置路径修改
import { get<%==s options.ApiName%>List, delete<%==s options.ApiName%> } from "@/api/index";
<%
}
%>
// 根据你的实际添加编辑组件位置路径修改
import add from "./add.vue";
import edit from "./edit.vue";
<%
if options.SearchDict{
%>
import dict from "@/views/my-components/xboot/dict";
<%
}
%>
<%
if options.SearchCustomList{
%>
import customList from "@/views/my-components/xboot/custom-list";
<%
}
%>
import { shortcuts } from "@/libs/shortcuts";
export default {
  name: "<%==s options.VueName%>",
  components: {
    add,
    edit,
    <%
    if options.SearchDict{
    %>
    dict,
    <%
    }
    %>
    <%
    if options.SearchCustomList{
    %>
    customList,
    <%
    }
    %>
  },
  data() {
    return {
      <% if options.SearchSize > 0 { %>
      openSearch: true, // 显示搜索
      <% } %>
      openTip: true, // 显示提示
      formData: {},
      currView: "index",
      loading: true, // 表单加载状态
      <% if options.HideSearch { %>
      drop: false,
      dropDownContent: "展开",
      dropDownIcon: "ios-arrow-down",
      <% } %>
      initSearchForm: {},
      searchForm: { // 搜索框初始化对象
        pageNumber: { value: 1 }, // 当前页数
        pageSize: { value: 10 }, // 页面大小
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
    submited() {
      this.currView = "index";
      this.getDataList();
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
    clearSelectAll() {
      this.$refs.table.selectAll(false);
    },
    changeSelect(e) {
      this.selectList = e;
    },
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
    add() {
      this.currView = "add";
    },
    edit(v) {
      // 转换null为""
      for (let attr in v) {
        if (v[attr] == null) {
          v[attr] = "";
        }
      }
      let str = JSON.stringify(v);
      let data = JSON.parse(str);
      this.formData = data;
      this.currView = "edit";
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