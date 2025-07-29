@echo off
chcp 65001 > nul
set command=%1

if "%command%" == "" (
	goto Help
) else if "%command%" == "tpl" (
	goto Gen_Tpl
) else (
	goto Help
)

:Gen_Tpl
@echo ◈ 更新VUE模板中...
hero -source="graph/extgen/template" -extensions=".tpl"
hero -source="template" -extensions=".tpl"
@echo ✓ 数据模型CRUD、VUE生成模板更新完成
goto End

:Help
@echo 参数错误,命令正确用法如下：
@echo core [command]
@echo   - command:
@echo     - tpl 更新数据模型CRUD、VUE生成模板文件,一般不常用
@echo     - 如提示hero 命令不存在，请先安装hero命令工具：
@echo     - go install github.com/shiyanhui/hero/hero
:End

