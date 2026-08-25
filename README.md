# kinship

家族图谱与亲缘关系推算 CLI。给定一个简单的家族登记文件，回答：
某人的祖先/子女有谁、任意两人之间是什么亲缘关系（mother、
great-uncle、half-sibling、first cousin once removed …）。

纯 Go 标准库实现，离线可构建。

## 构建

```sh
go build ./...
```

## 登记文件格式

每行一条记录，空行与 `#` 开头的行为注释：

```
PERSON <姓名> <F|M> <出生年>
PARENT <父/母姓名> <子女姓名>
```

- 姓名不能重复；PARENT 引用的姓名必须存在（允许前向引用）。
- 不允许自为父母、不允许重复的 PARENT 对。
- 语义错误均带行号报告。

## 子命令

| 命令 | 说明 |
| --- | --- |
| `kinship list <family.txt>` | 列出登记在册的所有人 |
| `kinship ancestors <family.txt> <name>` | 按辈分分组列出祖先 |
| `kinship children <family.txt> <name>` | 列出子女 |
| `kinship kin <family.txt> <a> <b>` | 输出 b 相对 a 的亲缘称谓 |

## 示例

```sh
kinship ancestors example/family.txt grace
# grandparents: carl, diana
# great-grandparents: george

kinship kin example/family.txt grace ivy
# grace -> ivy: first cousin once removed
```

## 亲缘规则

- 直系：parent/grandparent/great-grandparent，daughter/son 等，按性别选词。
- 兄弟姐妹：共享 2 名登记父母为 sibling，共享 1 名为 half-sibling。
- 旁系：aunt/uncle（含 great- 前缀）、niece/nephew。
- 表亲：以最低共同祖先距离 da、db 计算，min(da,db)-1 为第几代
  cousin，|da-db| 为 removed 层数。

## 测试

```sh
go test ./...
```
