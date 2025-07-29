<%!
import (
  "github.com/anden007/dp_clean_core/pkg"
)

type AddOption struct{
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
<%: func Add(fields []pkg.FormField,options AddOption, buffer *bytes.Buffer) %>
<template>
  <div>
    <Card>
      <div slot="title">
        <div class="edit-head">
          <a @click="close" class="back-title">
            <Icon type="ios-arrow-back" />返回
          </a>
          <div class="head-name">添加</div>
          <span></span>
          <a @click="close" class="window-close">
            <Icon type="ios-close" size="31" class="ivu-icon-ios-close" />
          </a>
        </div>
      </div>
      <div style="width: <%==s options.ModalWidth%>px">
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
          <FormItem class="br">
            <Button @click="handleSubmit" :loading="submitLoading" type="primary">提交并保存</Button>
            <Button @click="handleReset">重置</Button>
            <Button type="dashed" @click="close">关闭</Button>
          </FormItem>
        </Form>
      </div>
    </Card>
  </div>
</template>

<script>
<%
if options.Api{
%>
// 根据你的实际请求api.js位置路径修改
import { add<%==s options.ApiName%> } from "./api";
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
  name: "add",
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
      submitLoading: false, // 表单提交状态
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
    handleReset() {
      this.$refs.form.resetFields();
    },
    handleSubmit() {
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
          this.submitLoading = true;
          <%
          if options.Api{
          %>
          delete this.form.id;
          if('createTime' in this.form){
            delete this.form.createTime;
          }
          add<%==s options.ApiName%>(this.form).then(res => {
            this.submitLoading = false;
            if (res.success) {
              this.$Message.success("操作成功");
              this.submited();
            }
          });
          <%
          } else {
          %>
          // this.postRequest("请求路径", this.form).then(res => {
          //   this.submitLoading = false;
          //   if (res.success) {
          //     this.$Message.success("添加成功");
          //     this.submited();
          //   }
          // });
          // 模拟成功
          this.submitLoading = false;
          this.$Message.success("添加成功");
          this.submited();
          <%
          }
          %>
        }
      });
    },
    close() {
      this.$emit("close", true);
    },
    submited() {
      this.$emit("submited", true);
    }
  },
  mounted() {
    this.init();
  }
};
</script>
<style lang="less">
// 建议引入通用样式 具体路径自行修改 可删除下面样式代码
// @import "@/styles/single-common.less";
.edit-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    position: relative;

    .back-title {
        color: #515a6e;
        display: flex;
        align-items: center;
    }

    .head-name {
        display: inline-block;
        height: 20px;
        line-height: 20px;
        font-size: 16px;
        color: #17233d;
        font-weight: 500;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .window-close {
        z-index: 1;
        font-size: 12px;
        position: absolute;
        right: 0px;
        top: -5px;
        overflow: hidden;
        cursor: pointer;

        .ivu-icon-ios-close {
            color: #999;
            transition: color .2s ease;
        }
    }

    .window-close .ivu-icon-ios-close:hover {
        color: #444;
    }
}
</style>