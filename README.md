# Go 言語 + OpenGL + GLFW チュートリアル

このリポジトリは、" [○○くんのために一所懸命書いたものの
結局○○くんの卒業に間に合わなかった
GLFW による OpenGL 入門](http://marina.sys.wakayama-u.ac.jp/~tokoi/GLFWdraft.pdf) "の Go 言語による実装です。

OpenGL と Go 言語の練習用コードなので、参考までにどうぞ。

## 依存関係

以下の OpenGL 及び GLFW のラッパライブラリを使用しています。 `go get` コマンドでインストールして使用します。


```
github.com/go-gl/gl/v4.1-core/gl
github.com/go-gl/glfw/v3.2/glfw
```

OS は macOS 10.13 、Go 1.10.1 で動作確認を行っています。

## Usage

テキストの各チャプタに対応するコードは `chapter_n` のリポジトリの中に入っています。例えば 4 章のサンプルコードは、

```
$ go run chapter_4/main.go
```

のように実行してください。

## 参考リンク集

- [○○くんのために一所懸命書いたものの
結局○○くんの卒業に間に合わなかった
GLFW による OpenGL 入門](http://marina.sys.wakayama-u.ac.jp/~tokoi/GLFWdraft.pdf)
- [go-gl/example](https://github.com/go-gl/example)
  - glfw 及び gl ライブラリのサンプルコード
