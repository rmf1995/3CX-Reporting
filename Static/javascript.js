$.ajax({
    url: 'https://URL.com/v1/getAll3CX/',
    type: "get",
    dataType: "json",

    success: function(data) {
        drawTable(data);
    }
});

function drawTable(data) {
    for (var i = 0; i < data.length; i++) {
        drawRow(data[i]);
    }

    var table = $('#myTable').DataTable({
        fixedHeader: true,
        dom: 'Bfrtip',
        lengthMenu: [
            [10, 25, 50, -1],
            ['10 rows', '25 rows', '50 rows', 'Show all']
        ],
        buttons: ['pageLength', {
                extend: 'excelHtml5',
                title: '3CX_Report',
                exportOptions: {
                    columns: [0, ':visible']
                }
            },
            {
                extend: 'csvHtml5',
                title: '3CX_Report',
                exportOptions: {
                    columns: ':visible'
                }
            },
            {
                extend: 'pdfHtml5',
                title: '3CX_Report',
                exportOptions: {
                    columns: ':visible'
                }
            },
            'colvis'
        ]

    });




}

function epochToDate(timestamp) {
    var date = new Date(timestamp * 1000);
    var year = date.getFullYear();
    var month = date.getMonth() + 1;
    var day = date.getDate();
    var hours = date.getHours();
    var minutes = date.getMinutes();
    var seconds = date.getSeconds();
    return (year + "-" + month + "-" + day + " " + hours + ":" + minutes + ":" + seconds);
}


function inSpec(MaxSimCalls, vcpus, OSram) {
    if (MaxSimCalls >= 256) {
        if (vcpus == 8 && OSram == 16) {
            return true
        } else {
            return false
        }
    } else if (MaxSimCalls >= 96) {
        if (vcpus == 8 && OSram == 8) {
            return true
        } else {
            return false
        }
    } else if (MaxSimCalls >= 64) {
        if (vcpus == 6 && OSram == 6) {
            return true
        } else {
            return false
        }
    } else if (MaxSimCalls >= 48) {
        if (vcpus == 4 && OSram == 4) {
            return true
        } else {
            return false
        }
    } else if (MaxSimCalls >= 24) {
        if (vcpus == 4 && OSram == 4) {
            return true
        } else {
            return false
        }
    } else if (MaxSimCalls >= 16) {
        if (vcpus == 2 && OSram == 2) {
            return true
        } else {
            return false
        }
    } else if (MaxSimCalls >= 8) {
        if (vcpus == 2 && OSram == 2) {
            return true
        } else {
            return false
        }
    } else {
        return false
    }
}



function drawRow(rowData) {
    var row = $("<tr />")
    $("#myTable").append(row);
    row.append($("<td><a href=\"https://"+ rowData.URL +":5001\" target=\"_blank\">" + rowData.Name + "</a></td>"));
    row.append($("<td>" + rowData.CustomerID + "</td>"));
    row.append($("<td>" + rowData.Location + "</td>"));
    row.append($("<td>" + rowData.Version + "</td>"));
    row.append($("<td><a href=\"https://"+ rowData.FQDN +":5001\" target=\"_blank\">" + rowData.FQDN + "</a></td>"));
    row.append($("<td style=\"text-align:center\">" + rowData.CallRecordingUsage + "</td>"));
    row.append($("<td style=\"text-align:center\">" + rowData.MaxSimCalls + "</td>"));
    row.append($("<td style=\"text-align:center\">" + rowData.ExtTotal + "</td>"));
    row.append($("<td style=\"text-align:center\">" + rowData.vcpus + "</td>"));
    row.append($("<td style=\"text-align:center\">" + rowData.OSram + "</td>"));
    row.append($("<td style=\"text-align:center\">" + rowData.OSswap + "</td>"));
    row.append($("<td style=\"text-align:center\">" + rowData.OSDiskSpace + "</td>"));
    if (rowData.AutoUpdate == 0) {
        row.append($("<td style=\"text-align:center\">FALSE</td>"));
    } else {
        row.append($("<td style=\"text-align:center;background-color:red\">TRUE</td>"));
    }
    row.append($("<td style=\"text-align:center\">" + rowData.License + "</td>"));
    row.append($("<td style=\"text-align:center\">" + rowData.LicenseKey + "</td>"));
    if (rowData.LicenseExpiration == null) {
        row.append($("<td style=\"text-align:center;background-color:red\">" + epochToDate(rowData.LicenseExpiration) + "</td>"));
    } else if ((rowData.LicenseExpiration - (Date.now() / 1000)) <= 2678400) {
        row.append($("<td style=\"text-align:center;background-color:red\">" + epochToDate(rowData.LicenseExpiration) + "</td>"));
    } else {
        row.append($("<td style=\"text-align:center\">" + epochToDate(rowData.LicenseExpiration) + "</td>"));
    }
    row.append($("<td>" + rowData.ResellerName + "</td>"));
    if (rowData.lastUpdated == null) {
        row.append($("<td style=\"text-align:center;background-color:red\">" + epochToDate(rowData.lastUpdated) + "</td>"));
    } else if ((rowData.lastUpdated - (Date.now() / 1000)) >= 16200) {
        row.append($("<td style=\"text-align:center;background-color:red\">" + epochToDate(rowData.lastUpdated) + "</td>"));
    } else {
        row.append($("<td style=\"text-align:center\">" + epochToDate(rowData.lastUpdated) + "</td>"));
    }
    if (inSpec(rowData.MaxSimCalls, rowData.vcpus, rowData.OSram) == true) {
        row.append($("<td style=\"text-align:center;background-color:green\">" + "OK" + "</td>"));
    } else {
        row.append($("<td style=\"text-align:center;background-color:red\">" + "NOK" + "</td>"));
    }
}