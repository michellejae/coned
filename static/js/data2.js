const myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json");


const requestOptions = {
  method: "POST",
  headers: myHeaders,
};

async function foobar(){
  let graphData;

  try {
    const response = await fetch("/", requestOptions)
    graphData = await response.json()
  } catch (error) {
    console.log(error)
  }



var myChart = echarts.init(document.getElementById('main'));



// Specify the configuration items and data for the chart
var option = {
  tooltip: {
    formatter: function (params) {
      let index = params.dataIndex
      return `Name: ${params.name}<br />
      Total: $${params.value}<br />
      Rate: ${(Number(graphData[index].rate) * 100).toFixed(2)} ¢/kWh<br />
      Offer Type: ${graphData[index].offerType}<br />
      Minimum Contract Length: ${graphData[index].minTerm} months<br />
      Energy Source: ${graphData[index].energySource}<br />
      % Renewable: ${Number(graphData[index].percentRenew) * 100}% 
      `
    }
  },
  legend: {
  },
  xAxis: {
    type: "category",
    data: [],
    show: false
  },
  yAxis: {},
  series: [
    {
      type: 'bar',
      data: [],
    }
  ]
};



for (let i=0 ; i<graphData.length; i++){

  let name = graphData[i].name
  let total = graphData[i].total
 
  // setting all the properties that will get passed into data array 
  // basically each esco's is getting it's own obj
  let seriesDataObj = {}
  seriesDataObj.name = name
  seriesDataObj.value = total
  seriesDataObj.itemStyle = {}
  seriesDataObj.itemStyle.color = "blue"

  seriesDataObj.name === "Consolidated Edison Company of New York, Inc." ? seriesDataObj.itemStyle.color = "green" : seriesDataObj.itemStyle.color = "blue"
  
  //option.tooltip.formatter = `${name} <br/> ${total}`
  
  option.series[0].data.push(seriesDataObj)
}




// Display the chart using the configuration items and data just specified.
myChart.setOption(option);

}

foobar()


