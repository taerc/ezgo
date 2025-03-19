package main

func main() {
	// 生成含标题样式的文档
doc := document.New()
para := doc.AddParagraph()
para.SetStyle("Heading1")
para.AddRun().AddText("第一章 简介")
doc.SaveToFile("temp.docx")

// 调用 LibreOffice 命令行插入目录
cmd := exec.Command(
    "libreoffice", 
    "--headless", 
    "--invisible", 
    "--convert-to", "docx", 
    "--outdir", "/tmp", 
    "temp.docx", 
    "macro:///Standard.Module1.InsertTOC()" // 需提前录制宏
)
_ = cmd.Run()
}