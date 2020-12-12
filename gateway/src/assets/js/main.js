console.log("...Commerce Data Gateway!")

function csv(url, divObj){
    (async () => {
        try {
          console.log('CSV loading...');
          var response = await fetch(url);
          var data = await response.text();
          console.log(data);
          successFunction(data, divObj)
        } catch (e) {
          console.log('Error:', e);
          successFunction("error,"+e, divObj)
        }
      })();
}