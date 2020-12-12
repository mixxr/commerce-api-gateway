function successFunction(data, docDiv) {
    var allRows = data.split(/\r?\n|\r/);
    var table = '';
    for (var singleRow = 0; singleRow < allRows.length; singleRow++) {
      
        table += '<tr>';
      
      var rowCells = allRows[singleRow].split(',');
      for (var rowCell = 0; rowCell < rowCells.length; rowCell++) {

          table += '<td>';
          table += rowCells[rowCell];
          table += '</td>';

      }

        table += '</tr>';

    } 

    
    docDiv.innerHTML = table
  }