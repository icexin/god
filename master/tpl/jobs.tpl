<html>
<body>
<table border=1>
<tr>
<th>id</th>
<th>start_time</th>
<th>end_time</th>
<th>status</th>
<th>detail</th>
</tr>
{{range .}}
<tr>
<td>{{.Id}}</td>
<td>{{.StartTime}}</td>
<td>{{.EndTime}}</td>
<td>{{.Status}}</td>
<td><a href="/job/{{.Id}}">detail</a></td>
</tr>
{{end}}
</table>
</body>
</html>