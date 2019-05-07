<html>
<head>
	<title>上传文件</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
  <input type="file" name="uploadfile" />
  <input type="hidden" name="token" value="{{.token}}"/>
  <input type="submit" value="upload" />
</form>
姓名：{{ .user.UserName }}<br/>
Email：{{range .user.Emails}}
			{{.|html}}
		{{end}}<br/>
		
朋友：{{with .user.Friends}}
		{{range .}}
		{{.Fname}}
		{{end}}
		{{end}}<br/>
{{if `aaa`}}if有输出{{end}}<br/>
{{if ``}}if没有输出{{end}}<br/>
{{if ``}}if条件{{else}}else条件{{end}}<br/>
</body>
</html>