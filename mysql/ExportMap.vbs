Class Excel
  Private xl

  Private Sub Class_Initialize
    Set xl = CreateObject("Excel.Application")
  End Sub

  Private Sub Class_Terminate
    For Each wb In xl.Workbooks
      wb.Saved = True  'discard unsaved changes
      wb.Close         'close workbook
    Next
	xl.DisplayAlerts = False
    xl.Quit            'quit Excel
	Set xl = Nothing
  End Sub

  Public Function OpenWorkbook(filename)
    Set OpenWorkbook = xl.Workbooks.Open(filename)
  End Function

  Public Function RunMacro()
    Dim sht
    Dim cht
    Dim i
    i = 1

    ' xl.Run  "Module1.ExportCharts"
	xl.Application.Visible = False
	xl.Application.DisplayAlerts = False
    xl.Application.ScreenUpdating = False
    xl.Application.EnableEvents = False

    Set sht = xl.ActiveSheet

    For Each cht In sht.ChartObjects
        cht.Activate
    xl.ActiveChart.Export xl.ActiveWorkbook.Path & "\images\" & "graph-Map-" & i & ".png"
    i = i + 1
    Next

    xl.Application.EnableEvents = True
    xl.Application.ScreenUpdating = True
  End Function
  
  Public Function NewWorkbook
    Set NewWorkbook = xl.Workbooks.Add
  End Function

  Public Property Get Workbooks
    Set Workbooks = xl.Workbooks
  End Property
End Class

Dim Exl
Dim Arg
Dim vntCurrentDirectory
Dim objWshShell
Dim strExcelFileName
Dim wb1

Set objWshShell = WScript.CreateObject("WScript.Shell")
Set objFSO = CreateObject("Scripting.FileSystemObject")
vntCurrentDirectory = objWshShell.CurrentDirectory

Set Arg = WScript.Arguments
if WScript.Arguments.Count < 1 then
	WScript.Echo "Missing parameters"
end if

strExcelFileName   = Wscript.Arguments.Item(0)

Set xl = New Excel
Set wb1 = xl.OpenWorkbook(vntCurrentDirectory & strExcelFileName)
xl.RunMacro()