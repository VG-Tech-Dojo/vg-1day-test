package api

import (
	"model"
	"request"

	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/labstack/echo"
)

// メッセージの投稿
func CreateMessage(c echo.Context) error {
	// request.Message を用意する
	var r request.Message

	// 受け取った json を request.Message として取得する
	if err := c.Bind(&r); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return err
	}

	// メッセージをつくる
	message, err := model.NewMessage(r.Body, r.Username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return err
	}

	// メッセージを保存する
	if err := message.SaveMessage(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return err
	}

	// メッセージを json で返す
	return c.JSON(http.StatusCreated, message)
}

// メッセージの取得
func ReadMessage(c echo.Context) error {
	// id を受け取る
	id, _ := strconv.Atoi(c.Param("id"))

	// model.Message を用意する
	var m model.Message

	// メッセージを取得する
	if err := m.LoadMessage(id); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return err
	} else {
		// メッセージを json で返す
		return c.JSON(http.StatusOK, m)
	}
}

// 1-3. メッセージの更新
func UpdateMessage(c echo.Context) error {
	var m model.Message

	if err := c.Bind(&m); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return err
	}

	id := m.Id
	body := m.Body

	if err := m.LoadMessage(id); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return err
	}

	m.Body = body

	if err := m.SaveMessage(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, m)
	return nil
}

// 1-4. メッセージの削除
func DeleteMessage(c echo.Context) error {
	// model.Message を用意する
	// 受け取った id を使って model.Message を取得する

	// メッセージを削除する
	// ヒント: model.Message.DeleteMessage()

	return c.NoContent(http.StatusOK)
}

// メッセージ一覧の取得
func ReadMessages(c echo.Context) error {
	// メッセージ一覧を取得する
	messages, err := model.LoadMessages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return err
	}

	if len(messages) == 0 {
		// 空の時は、空っぽで返す
		return c.NoContent(http.StatusNoContent)
	}

	// メッセージ一覧を json で返す
	return c.JSON(http.StatusOK, messages)
}

// メッセージの投稿 (bot 対応)
func ObservableCreateMessage(ch chan model.Message) echo.HandlerFunc {
	return func(c echo.Context) error {
		// request.Message を用意する
		var r request.Message

		// 受け取った json を request.Message として取得する
		if err := c.Bind(&r); err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			return err
		}

		// メッセージをつくる
		message, err := model.NewMessage(r.Body, r.Username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			return err
		}

		// メッセージを保存する
		if err := message.SaveMessage(); err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			return err
		}

		if ch != nil {
			ch <- *message
		}

		// メッセージを json で返す
		return c.JSON(http.StatusCreated, message)
	}
}

func Mutex(m *sync.Mutex) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m.Lock()
			defer m.Unlock()
			return next(c)
		}
	}
}
