module github.com/stratsys/plompose

go 1.13

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.6.0

replace github.com/docker/libcompose => github.com/docker/libcompose v0.4.1-0.20171025083809-57bd716502dc

replace github.com/docker/cli => github.com/docker/cli v0.0.0-20180529093712-df6e38b81a94

replace github.com/xeipuuv/gojsonschema => github.com/xeipuuv/gojsonschema v0.0.0-20160323030313-93e72a773fad

replace github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20180220021536-8e435b8279f2+incompatible

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20190830141801-acfa387b8d69

replace github.com/kubernetes/kompose => github.com/stratsys/kompose v1.22.2

require (
	github.com/docker/cli v0.0.0-00010101000000-000000000000
	github.com/kubernetes/kompose v1.22.0
	golang.org/x/sys v0.0.0-20200831180312-196b9ba8737a // indirect
)
