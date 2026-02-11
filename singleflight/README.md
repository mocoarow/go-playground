# singleflight

golang.org/x/sync/singleflight は、同じキーに対する重複した関数呼び出しを抑制する仕組みです。

## 問題

例えば、10個の goroutine が同時に同じユーザーデータを取得しようとすると、通常は 10 回 DB や API を叩きます

goroutine 1 → Fetch("user:123") → DB 呼び出し  
goroutine 2 → Fetch("user:123") → DB 呼び出し  ← 無駄  
goroutine 3 → Fetch("user:123") → DB 呼び出し  ← 無駄  
...

## singleflight の解決策

group.Do(key, fn) を使うと、同じキーで同時に呼ばれた場合、最初の1つだけが fn を実行し、残りはその結果を共有して待ちます。

goroutine 1 → Do("user:123", fn) → 実際に fn 実行 → 結果を返す  
goroutine 2 → Do("user:123", fn) → 待機...       → 同じ結果を受け取る  
goroutine 3 → Do("user:123", fn) → 待機...       → 同じ結果を受け取る  

## 主な API

| メソッド | 用途 |
|---------|------|
| `Do(key, fn)` | 重複抑制して実行。同じキーの呼び出しは結果を共有 |
| `DoChan(key, fn)` | `Do` の非同期版。チャネルで結果を受け取る |
| `Forget(key)` | キーの進行中エントリを削除。次の呼び出しで再実行される |

## 戻り値

```go
v, err, shared := group.Do(key, fn)
//      ^^^^^^ 結果が他の呼び出しと共有されたか
```

shared が true なら、他の goroutine と結果を共有したことを意味します。  
shared が false なら、その結果を受け取ったのは自分だけだったことを意味します。  

## 典型的なユースケース

- キャッシュの thundering herd 対策 - キャッシュ期限切れ時に大量リクエストがDBに殺到するのを防ぐ
- API ゲートウェイ - 同一リクエストの重複を排除
- DNS リゾルバ - 同一ホスト名の並行解決を1回にまとめる

## 注意点

- キャッシュではない - 実行完了後にキーは消える。次の呼び出しでは再実行される
- エラーも共有される - 1つの実行が失敗すると、待機中の全呼び出しにエラーが返る
- タイムアウト制御は自分で行う - context によるキャンセルは Do 自体にはないので、fn 内で対応する
