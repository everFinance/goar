package data

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/everFinance/goar"
	"github.com/stretchr/testify/assert"
)

func TestDecodeDataItem(t *testing.T) {
	t.Run("DecodeDataItem - New empty test data item", func(t *testing.T) {
		data := ""
		tags := []Tag{}
		anchor := ""
		target := ""
		s, err := goar.NewSignerFromPath("../data/wallet.json")
		assert.NoError(t, err)

		a, err := NewDataItem([]byte(data), s, target, anchor, tags)
		assert.NoError(t, err)

		dataItem, err := DecodeDataItem(a.Raw)
		assert.NoError(t, err)
		assert.Equal(t, dataItem.Owner, "gcUyzZCfsHfIk1WeHt2iMG5VktieAhakX8--cPa_Hi-Z0HICCe35ihgaGZXvkwIRI1oFvbC0DJmm5q3WUyx8NVa7PA1kxTirSjgaqw82fZjR_z6YL51Vki_REuTnkbJP1znjYR6ie5a1THppRGzwpPZK-uh6rQkMUOWZO5qHx_Jv3gAuEKu_IBH39Aef1BwPU3np_0vKkKQVMhpgyo7gDIxB2VYqL3WwAxRGuTf9x2Ihp3AU_dJPTJ6AuaLUR8b39YpYRe8bCWRjbOPlU4IL2_WfTGPUnIxnGXUzUMNZNjXy65zKhW3DrJgv48tk78_iVSWuX73EZPaJz521f7yuCINs84QnE0Q1-EPS5aX4yH4bHTGgHvXarDxORazR9wCEGvXLdYyTlt1DAmeHZSpdk5lNRVtBnbSHRYqQ0qJ02HPJn6WWDuNgtxxoEUEXSFojNQ5NGsDbxxjVEeo7cIdJXaF83E9p84Tvayc0q69Vu0pA8fE3ZL79mCT153nlnRHeAdUNX9H4vKXiGkomEO8Gun4Fi5dKfgOORtD5u0yMyeF__S1nT5a3bv29CLDbcgX5iWufSO9uXmF1atTo7NbcOKu-sm_tWwr4T95stF78dZu0DkzUO8ylGXF5r5Zzi2SVkOG1cjwLrZ91cAuPTNuzFfy04Sdv05fKoYZn-AJatvU")
		assert.Equal(t, dataItem.Target, target)
		assert.Equal(t, dataItem.Anchor, anchor)
		assert.Equal(t, dataItem.Data, base64.RawURLEncoding.EncodeToString([]byte(data)))
	})

	t.Run("DecodeDataItem - data, tags, anchor, target", func(t *testing.T) {
		data := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{};':\",./<>?`~"
		tags := []Tag{
			{Name: "tag1", Value: "value1"},
			{Name: "tag2", Value: "value2"},
		}
		anchor := "thisSentenceIs32BytesLongTrustMe"
		target := "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs"

		s, err := goar.NewSignerFromPath("../data/wallet.json")
		assert.NoError(t, err)

		a, err := NewDataItem([]byte(data), s, target, anchor, tags)
		assert.NoError(t, err)

		dataItem, err := DecodeDataItem(a.Raw)
		assert.NoError(t, err)

		assert.Equal(t, dataItem.Owner, "gcUyzZCfsHfIk1WeHt2iMG5VktieAhakX8--cPa_Hi-Z0HICCe35ihgaGZXvkwIRI1oFvbC0DJmm5q3WUyx8NVa7PA1kxTirSjgaqw82fZjR_z6YL51Vki_REuTnkbJP1znjYR6ie5a1THppRGzwpPZK-uh6rQkMUOWZO5qHx_Jv3gAuEKu_IBH39Aef1BwPU3np_0vKkKQVMhpgyo7gDIxB2VYqL3WwAxRGuTf9x2Ihp3AU_dJPTJ6AuaLUR8b39YpYRe8bCWRjbOPlU4IL2_WfTGPUnIxnGXUzUMNZNjXy65zKhW3DrJgv48tk78_iVSWuX73EZPaJz521f7yuCINs84QnE0Q1-EPS5aX4yH4bHTGgHvXarDxORazR9wCEGvXLdYyTlt1DAmeHZSpdk5lNRVtBnbSHRYqQ0qJ02HPJn6WWDuNgtxxoEUEXSFojNQ5NGsDbxxjVEeo7cIdJXaF83E9p84Tvayc0q69Vu0pA8fE3ZL79mCT153nlnRHeAdUNX9H4vKXiGkomEO8Gun4Fi5dKfgOORtD5u0yMyeF__S1nT5a3bv29CLDbcgX5iWufSO9uXmF1atTo7NbcOKu-sm_tWwr4T95stF78dZu0DkzUO8ylGXF5r5Zzi2SVkOG1cjwLrZ91cAuPTNuzFfy04Sdv05fKoYZn-AJatvU")
		assert.Equal(t, dataItem.Target, target)
		assert.Equal(t, dataItem.Anchor, anchor)
		assert.Equal(t, dataItem.Data, base64.RawURLEncoding.EncodeToString([]byte(data)))
	})
	t.Run("DecodeDataItem - Stub", func(t *testing.T) {
		data, err := os.ReadFile("../test/stubs/1115BDataItem")
		assert.NoError(t, err)

		dataItem, err := DecodeDataItem(data)
		assert.NoError(t, err)
		assert.Equal(t, dataItem.ID, "QpmY8mZmFEC8RxNsgbxSV6e36OF6quIYaPRKzvUco0o")
		assert.Equal(t, dataItem.Signature, "wUIlPaBflf54QyfiCkLnQcfakgcS5B4Pld-hlOJKyALY82xpAivoc0fxBJWjoeg3zy9aXz8WwCs_0t0MaepMBz2bQljRrVXnsyWUN-CYYfKv0RRglOl-kCmTiy45Ox13LPMATeJADFqkBoQKnGhyyxW81YfuPnVlogFWSz1XHQgHxrFMAeTe9epvBK8OCnYqDjch4pwyYUFrk48JFjHM3-I2kcQnm2dAFzFTfO-nnkdQ7ulP3eoAUr-W-KAGtPfWdJKFFgWFCkr_FuNyHYQScQo-FVOwIsvj_PVWEU179NwiqfkZtnN8VoBgCSxbL1Wmh4NYL-GsRbKz_94hpcj5RiIgq0_H5dzAp-bIb49M4SP-DcuIJ5oT2v2AfPWvznokDDVTeikQJxCD2n9usBOJRpLw_P724Yurbl30eNow0U-Jmrl8S6N64cjwKVLI-hBUfcpviksKEF5_I4XCyciW0TvZj1GxK6ET9lx0s6jFMBf27-GrFx6ZDJUBncX6w8nDvuL6A8TG_ILGNQU_EDoW7iil6NcHn5w11yS_yLkqG6dw_zuC1Vkg1tbcKY3703tmbF-jMEZUvJ6oN8vRwwodinJjzGdj7bxmkUPThwVWedCc8wCR3Ak4OkIGASLMUahSiOkYmELbmwq5II-1Txp2gDPjCpAf9gT6Iu0heAaXhjk")
		assert.Equal(t, dataItem.Owner, "0zBGbs8Y4wvdS58cAVyxp7mDffScOkbjh50ZrqnWKR_5NGwjezT6J40ejIg5cm1KnuDnw9OhvA7zO6sv1hEE6IaGNnNJWiXFecRMxCl7iw78frrT8xJvhBgtD4fBCV7eIvydqLoMl8K47sacTUxEGseaLfUdYVJ5CSock5SktEEdqqoe3MAso7x4ZsB5CGrbumNcCTifr2mMsrBytocSoHuiCEi7-Nwv4CqzB6oqymBtEECmKYWdINnNQHVyKK1l0XP1hzByHv_WmhouTPos9Y77sgewZrvLF-dGPNWSc6LaYGy5IphCnq9ACFrEbwkiCRgZHnKsRFH0dfGaCgGb3GZE-uspmICJokJ9CwDPDJoxkCBEF0tcLSIA9_ofiJXaZXbrZzu3TUXWU3LQiTqYr4j5gj_7uTclewbyZSsY-msfbFQlaACc02nQkEkr4pMdpEOdAXjWP6qu7AJqoBPNtDPBqWbdfsLXgyK90NbYmf3x4giAmk8L9REy7SGYugG4VyqG39pNQy_hdpXdcfyE0ftCr5tSHVpMreJ0ni7v3IDCbjZFcvcHp0H6f6WPfNCoHg1BM6rHUqkXWd84gdHUzo9LTGq9-7wSBCizpcc_12_I-6yvZsROJvdfYOmjPnd5llefa_X3X1dVm5FPYFIabydGlh1Vs656rRu4dzeEQwc")
		assert.Equal(t, dataItem.Target, "")
		assert.Equal(t, dataItem.Anchor, "")
		assert.ElementsMatch(
			t,
			dataItem.Tags,
			[]Tag{
				{Name: "Content-Type", Value: "text/plain"},
				{Name: "App-Name", Value: "ArDrive-CLI"},
				{Name: "App-Version", Value: "1.21.0"},
			},
		)
		assert.Equal(t, dataItem.Data, "NTY3MAo")
	})
}

func TestNewDataItem(t *testing.T) {
	t.Run("NewDataItem - New empty test data item", func(t *testing.T) {
		data := ""
		tags := []Tag{}
		anchor := ""
		target := ""

		s, err := goar.NewSignerFromPath("../data/wallet.json")
		assert.NoError(t, err)
		dataItem, err := NewDataItem([]byte(data), s, target, anchor, tags)
		assert.NoError(t, err)

		assert.Equal(t, dataItem.Owner, "gcUyzZCfsHfIk1WeHt2iMG5VktieAhakX8--cPa_Hi-Z0HICCe35ihgaGZXvkwIRI1oFvbC0DJmm5q3WUyx8NVa7PA1kxTirSjgaqw82fZjR_z6YL51Vki_REuTnkbJP1znjYR6ie5a1THppRGzwpPZK-uh6rQkMUOWZO5qHx_Jv3gAuEKu_IBH39Aef1BwPU3np_0vKkKQVMhpgyo7gDIxB2VYqL3WwAxRGuTf9x2Ihp3AU_dJPTJ6AuaLUR8b39YpYRe8bCWRjbOPlU4IL2_WfTGPUnIxnGXUzUMNZNjXy65zKhW3DrJgv48tk78_iVSWuX73EZPaJz521f7yuCINs84QnE0Q1-EPS5aX4yH4bHTGgHvXarDxORazR9wCEGvXLdYyTlt1DAmeHZSpdk5lNRVtBnbSHRYqQ0qJ02HPJn6WWDuNgtxxoEUEXSFojNQ5NGsDbxxjVEeo7cIdJXaF83E9p84Tvayc0q69Vu0pA8fE3ZL79mCT153nlnRHeAdUNX9H4vKXiGkomEO8Gun4Fi5dKfgOORtD5u0yMyeF__S1nT5a3bv29CLDbcgX5iWufSO9uXmF1atTo7NbcOKu-sm_tWwr4T95stF78dZu0DkzUO8ylGXF5r5Zzi2SVkOG1cjwLrZ91cAuPTNuzFfy04Sdv05fKoYZn-AJatvU")
		assert.Equal(t, dataItem.Target, target)
		assert.Equal(t, dataItem.Anchor, anchor)
		assert.Equal(t, dataItem.Data, base64.RawURLEncoding.EncodeToString([]byte(data)))
	})
	t.Run("NewDataItem - data, tags, anchor, target", func(t *testing.T) {
		data := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{};':\",./<>?`~"
		tags := []Tag{
			{Name: "tag1", Value: "value1"},
			{Name: "tag2", Value: "value2"},
		}
		anchor := "thisSentenceIs32BytesLongTrustMe"
		target := "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs"

		s, err := goar.NewSignerFromPath("../data/wallet.json")
		assert.NoError(t, err)
		dataItem, err := NewDataItem([]byte(data), s, target, anchor, tags)
		assert.NoError(t, err)

		assert.Equal(t, dataItem.Owner, "gcUyzZCfsHfIk1WeHt2iMG5VktieAhakX8--cPa_Hi-Z0HICCe35ihgaGZXvkwIRI1oFvbC0DJmm5q3WUyx8NVa7PA1kxTirSjgaqw82fZjR_z6YL51Vki_REuTnkbJP1znjYR6ie5a1THppRGzwpPZK-uh6rQkMUOWZO5qHx_Jv3gAuEKu_IBH39Aef1BwPU3np_0vKkKQVMhpgyo7gDIxB2VYqL3WwAxRGuTf9x2Ihp3AU_dJPTJ6AuaLUR8b39YpYRe8bCWRjbOPlU4IL2_WfTGPUnIxnGXUzUMNZNjXy65zKhW3DrJgv48tk78_iVSWuX73EZPaJz521f7yuCINs84QnE0Q1-EPS5aX4yH4bHTGgHvXarDxORazR9wCEGvXLdYyTlt1DAmeHZSpdk5lNRVtBnbSHRYqQ0qJ02HPJn6WWDuNgtxxoEUEXSFojNQ5NGsDbxxjVEeo7cIdJXaF83E9p84Tvayc0q69Vu0pA8fE3ZL79mCT153nlnRHeAdUNX9H4vKXiGkomEO8Gun4Fi5dKfgOORtD5u0yMyeF__S1nT5a3bv29CLDbcgX5iWufSO9uXmF1atTo7NbcOKu-sm_tWwr4T95stF78dZu0DkzUO8ylGXF5r5Zzi2SVkOG1cjwLrZ91cAuPTNuzFfy04Sdv05fKoYZn-AJatvU")
		assert.Equal(t, dataItem.Target, target)
		assert.Equal(t, dataItem.Anchor, anchor)
		assert.Equal(t, dataItem.Data, base64.RawURLEncoding.EncodeToString([]byte(data)))
	})
}

func TestVerifyDataItem(t *testing.T) {
	t.Run("VerifyDataItem - Empty test data item", func(t *testing.T) {
		data := ""
		tags := []Tag{}
		anchor := ""
		target := ""

		s, err := goar.NewSignerFromPath("../data/wallet.json")
		assert.NoError(t, err)

		dataItem, err := NewDataItem([]byte(data), s, target, anchor, tags)
		assert.NoError(t, err)

		valid, err := VerifyDataItem(dataItem)
		assert.NoError(t, err)
		assert.True(t, valid)
	})
	t.Run("VerifyDataItem - data, tags, anchor, target", func(t *testing.T) {
		data := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{};':\",./<>?`~"
		tags := []Tag{
			{Name: "tag1", Value: "value1"},
			{Name: "tag2", Value: "value2"},
		}
		anchor := "thisSentenceIs32BytesLongTrustMe"
		target := "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs"

		s, err := goar.NewSignerFromPath("../data/wallet.json")
		assert.NoError(t, err)

		dataItem, err := NewDataItem([]byte(data), s, target, anchor, tags)
		assert.NoError(t, err)

		valid, err := VerifyDataItem(dataItem)
		assert.NoError(t, err)
		assert.True(t, valid)
	})
	t.Run("VerifyDataItem - Stub", func(t *testing.T) {
		data, err := os.ReadFile("../test/stubs/1115BDataItem")
		assert.NoError(t, err)

		dataItem, err := DecodeDataItem(data)
		assert.NoError(t, err)

		valid, err := VerifyDataItem(dataItem)
		assert.NoError(t, err)
		assert.True(t, valid)
	})
}
