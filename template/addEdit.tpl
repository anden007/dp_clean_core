<%!
import (
  "github.com/anden007/dp_clean_core/pkg"
)

type AddEditOption struct{
    Api bool
    ApiName string
    CustomList bool
    Dict bool
    Password bool
    Editor bool
    Upload bool
    UploadThumb bool
    FileUpload bool
    TotalRow int
    RowNum int
    Span string
    LabelPosition string
    ModalWidth string
}
%>
<%: func AddEdit(fields []pkg.FormField, options AddEditOption, buffer *bytes.Buffer) %>
<template>
  <div>
    <!-- Drawer抽屉 -->
    <Drawer :title="title" v-model="visible" width="<%==s options.ModalWidth%>" draggable :mask-closable="type=='0'">
      <div :style="{maxHeight: maxHeight}" class="drawer-content">
        <Form ref="form" :model="form" :rules="formValidate" :style="{minHeight: maxHeight}" <% if options.LabelPosition=="left"{ %>:label-width="100"<% } %> label-position="<%==s options.LabelPosition%>">
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
                <Select v-model="form.<%==s item.Field%>" clearable transfer>
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
                <DatePicker type="date" v-model="form.<%==s item.Field%>" clearable transfer style="width: 100%"></DatePicker>
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
                <DatePicker type="daterange" v-model="form.<%==s item.Field%>" clearable transfer style="width: 100%"></DatePicker>
                <%
                }
                %>
                <%
                if item.Type=="time"{
                %>
                <TimePicker type="time" v-model="form.<%==s item.Field%>" clearable transfer style="width: 100%"></TimePicker>
                <%
                }
                %>
                <%
                if item.Type=="area"{
                %>
                <al-cascader v-model="form.<%==s item.Field%>" data-type="name" level="<%==s item.Level%>" transfer/>
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
      </div>
      <div class="drawer-footer br" v-show="type!='0'">
        <Button type="primary" :loading="submitLoading" @click="submit">提交</Button>
        <Button @click="visible = false">取消</Button>
      </div>
    </Drawer>
  </div>
</template>

<script>
<%
if options.Api{
%>
// 根据你的实际请求api.js位置路径修改
import { add<%==s options.ApiName%>, edit<%==s options.ApiName%> } from "./api";
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
  name: "addEdit",
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
  props: {
    value: {
      type: Boolean,
      default: false
    },
    data: {
      type: Object
    },
    type: {
      type: String,
      default: "0"
    }
  },
  data() {
    return {
      visible: this.value,
      title: "",
      submitLoading: false, // 表单提交状态
      maxHeight: 510,
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
        for _, item := range fields{
          if item.Editable && item.Validate{
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
      }
    };
  },
  methods: {
    init() {},
    submit() {
      this.$refs.form.validate(valid => {
        if (valid) {
          <%
          for _,item := range fields{
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
          if (this.type == "1") {
            // 编辑
            this.submitLoading = true;
            <%
            if options.Api{
            %>
            edit<%==s options.ApiName%>(this.form).then(res => {
              this.submitLoading = false;
              if (res.success) {
                this.$Message.success("操作成功");
                this.$emit("on-submit", true);
                this.visible = false;
              }
            });
            <%
            } else {
            %>
            // this.postRequest("请求路径", this.form).then(res => {
            //   this.submitLoading = false;
            //   if (res.success) {
            //     this.$Message.success("操作成功");
            //     this.$emit("on-submit", true);
            //     this.visible = false;
            //   }
            // });
            // 模拟请求
            this.submitLoading = false;

            this.$Message.success("操作成功");
            this.$emit("on-submit", true);
            this.visible = false;
            <%
            }
            %>
          } else {
            // 添加
            this.submitLoading = true;
            <%
            if options.Api{
            %>
            delete this.form.id;
            add<%==s options.ApiName%>(this.form).then(res => {
              this.submitLoading = false;
              if (res.success) {
                this.$Message.success("操作成功");
                this.$emit("on-submit", true);
                this.visible = false;
              }
            });
            <%
            } else {
            %>
            // this.postRequest("请求路径", this.form).then(res => {
            //   this.submitLoading = false;
            //   if (res.success) {
            //     this.$Message.success("操作成功");
            //     this.$emit("on-submit", true);
            //     this.visible = false;
            //   }
            // });
            // 模拟请求
            this.submitLoading = false;
            this.$Message.success("操作成功");
            this.$emit("on-submit", true);
            this.visible = false;
            <%
            }
            %>
          }
        }
      });
    },
    setCurrentValue(value) {
      if (value === this.visible) {
        return;
      }
      if (this.type == "1") {
        this.title = "编辑";
        this.maxHeight =
          Number(document.documentElement.clientHeight - 121) + "px";
      } else if (this.type == "2") {
        this.title = "添加";
        this.maxHeight =
          Number(document.documentElement.clientHeight - 121) + "px";
      } else {
        this.title = "详细信息";
        this.maxHeight = "100%";
      }
      // 清空数据
      this.$refs.form.resetFields();
      if (this.type == "0" || this.type == "1") {
        // 回显数据处理
        this.form = this.data;
      } else {
        // 添加
        delete this.form.id;
      }
      this.visible = value;
    }
  },
  watch: {
    value(val) {
      this.setCurrentValue(val);
    },
    visible(value) {
      this.$emit("input", value);
    }
  },
  mounted() {
    this.init();
  }
};
</script>

<style lang="less">
// 建议引入通用样式 具体路径自行修改 可删除下面样式代码
// @import "@/styles/drawer-common.less";
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