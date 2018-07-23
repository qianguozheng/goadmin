package admin

import (
	"net/http"

	"github.com/labstack/echo"
)

//html template usage
//https://www.cnblogs.com/Pynix/p/4154630.html

type Test struct {
	Name    string
	Md5     string
	Version string
}

func HandleExample(c echo.Context) error {
	//Table 2 column
	pair := make(map[string]string)
	pair["hello"] = "jim"
	pair["jon"] = "snow"
	pair["Richard"] = "Qian"

	//Table 3 column
	triple := make(map[string]Test)
	test1 := Test{
		Name:    "Test1",
		Md5:     "md51",
		Version: "version1",
	}

	test2 := Test{
		Name:    "Test2",
		Md5:     "md53",
		Version: "versio4",
	}

	triple["9341"] = test1
	triple["9344"] = test2

	var slices []Test
	slices = append(slices, test1)
	slices = append(slices, test2)
	return c.Render(http.StatusOK, "example.html", map[string]interface{}{
		"Pair":   pair,
		"Triple": triple,
		"Slice":  slices,
	})
}
