<html>
<body>
<table border=1>
<tr>
<th>agent</th>
<th>start_time</th>
<th>end_time</th>
<th>exit_code</th>
<th>status</th>
<th>output</th>
</tr>
{{range .}}
<tr>
<td>{{.Agent}}</td>
<td>{{.StartTime}}</td>
<td>{{.EndTime}}</td>
<td>{{.ExitCode}}</td>
<td>{{.Status}}</td>
<td><a href="/job/{{.Id}}/{{.Agent}}/output">output</a></td>
</tr>
{{end}}
</table>
</body>
</html>
