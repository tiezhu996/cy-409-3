# 法规条文全文检索工具

法规条文全文检索工具是一个 Go CLI，用于导入法律法规 TXT/Markdown 文本，并按关键词、拼音、条款号和章节进行本地全文检索，数据存储在 BoltDB。

## 安装方式

```bash
go mod download
go build -o lawsearch .
go install .
```

## 使用示例

```bash
lawsearch import --file 民法典.txt --name "中华人民共和国民法典"
lawsearch search --query "合同解除" --law "中华人民共和国民法典"
lawsearch search --query "hejie" --fuzzy --law "中华人民共和国民法典"
lawsearch search --query "侵权责任" --law "中华人民共和国民法典" --output markdown --file result.md
lawsearch article --law "中华人民共和国民法典" --number 1042
lawsearch article --law "中华人民共和国民法典" --from 100 --to 200
lawsearch toc --law "中华人民共和国民法典"
lawsearch toc --law "中华人民共和国民法典" --chapter "第三编 合同"
lawsearch list
lawsearch info --law "中华人民共和国民法典"
lawsearch delete --law "中华人民共和国民法典"
```

所有命令都支持 `--db lawsearch.db` 指定 BoltDB 数据库文件。

## 功能列表

- 法规导入：解析编、章、节、条结构，生成法规元信息和条款列表。
- 精确检索：支持单关键词和多关键词 AND/OR 组合。
- 模糊检索：支持拼音匹配和编辑距离不超过 2 的容错匹配。
- 条款号查询：支持单条和范围查询。
- 章节浏览：以树形文本形式浏览法规目录。
- 结果导出：检索结果可导出 Markdown 或 JSON。
- 法规管理：列出、查看、删除已导入法规。

## 技术栈

| 模块 | 技术 |
| --- | --- |
| CLI | Cobra |
| 数据库 | BoltDB / bbolt |
| 中文分词 | gojieba |
| 拼音检索 | go-pinyin |
| 表格输出 | tablewriter |
| 语言 | Go 1.21+ |

## 目录结构

```text
cmd/
├── root.go
├── import.go
├── search.go
├── article.go
├── toc.go
├── list.go
├── info.go
└── delete.go
internal/
├── parser/
├── indexer/
├── search/
├── store/
└── export/
pkg/
└── models/
```

## License

MIT
