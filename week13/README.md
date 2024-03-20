# 第十二次作业

## 要求

在 Validator 里面，base -> target 的这个过程，我们都是一条条取出来比较的。

现在需要你修改为批量接口。

也就是，从 base 中取一批，而后从 target 里面找出对应的数据，比较是否相等。

**可选** 你可以测试一下，在 base 中有十万条数据的情况下，单个和批量校验所花时间的对比。

## 实现思路

1. fromBase 签名的返回值改成切片.
```go
fromBase func(ctx context.Context, offset int) ([]T, error)
```
2. fullFromBase & incrFromBase 实现改成查询一批.
3. validateBaseToTarget 将一批数据到目标库查询(使用`validateTargetBatch`方法).
4. validateTargetBatch 使用 `id in ?` 批量查询目标库.
5. 然后遍历原数据，判断目标数据有没有;如果有的话，看看是否相等。

### 快速跳转

1. [fromBase 签名的返回值改成切片](./webook/pkg/migrator/validator/validator.go#L29)
2. [fullFromBase & incrFromBase 实现改成查询一批.](./webook/pkg/migrator/validator/validator.go#L157)
3. [validateBaseToTarget 将一批数据到目标库查询(使用`validateTargetBatch`方法).](./webook/pkg/migrator/validator/validator.go#L61)
4. [validateTargetBatch实现](./webook/pkg/migrator/validator/validator.go#L91)

