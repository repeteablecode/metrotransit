const endpoint = "https://cors-anywhere.herokuapp.com/https://svc.metrotransit.org/nextrip"

let routeNumber = "0"
let directional = "0"
let stop = "XXXX"
let departuresList = []

let setStatus = (status) => document.querySelector("#status").innerText = status

let unhide = (className) => document.querySelector(`.${className}`).classList.remove('hidden')

let getRoutes = () => {
  let routesEndpoint = `${endpoint}/Routes?format=JSON`
  setStatus('fetching routes...')
  fetch(routesEndpoint)
  .then((response) => response.json())
  .then((routes) => {
    setStatus('routes received..')
    routes.forEach((route) => {
      let opt = document.createElement('option')
      let txt = document.createTextNode(route.Description)
      opt.value = route.Route
      opt.appendChild(txt)
      selRoute.appendChild(opt)
    })
    unhide('routes')
  })
  .catch((err) => { throw Error(err)} )
}

let routePicked = (event) => {
  let dex = event.target.selectedIndex
  routeNumber = event.target.options[dex].value
  let ep = `${endpoint}/Directions/${routeNumber}?format=JSON`
  setStatus('fetching directions...')
  fetch(ep)
  .then((response) => response.json())
  .then((directions) => {
    setStatus('directions received...')
    directions.forEach((direction) => {
      let opt = document.createElement('option')
      let txt = document.createTextNode(direction.Text)
      opt.value = direction.Value
      opt.appendChild(txt)
      selDirection.appendChild(opt)
    })
    unhide('directions')
  })
}

let directionPicked = (event) => {
  let dex = event.target.selectedIndex
  directional = event.target.options[dex].value
  let ep = `${endpoint}/Stops/${routeNumber}/${directional}?format=JSON`
  setStatus('fetching stops...')
  fetch(ep)
  .then((response) => response.json())
  .then((stops) => {
    setStatus('stops received...')
    stops.forEach((stop) => {
      let opt = document.createElement('option')
      let txt = document.createTextNode(stop.Text)
      opt.value = stop.Value
      opt.appendChild(txt)
      selStop.appendChild(opt)
    })
    unhide('stops')
  })
  
}

let stopPicked = (event) => {
  let dex = event.target.selectedIndex;
  stop = event.target.options[dex].value
  let ep = `${endpoint}/${routeNumber}/${directional}/${stop}?format=JSON`
  setStatus("fetching departures...")
  fetch(ep)
  .then((response) => response.json())
  .then((departures) => {
    setStatus("departures received...")
    departuresList = departures
    document.querySelector("#arrival").innerText = departures[0].DepartureText
    setStatus(`Bus arriving at ${departures[0].DepartureText}`)
    unhide('arrival')
  })

}

let loaded = () => {
  document.querySelector("#title").innerText = "Metro Transit NexTrip API"
  selRoute.addEventListener("change", routePicked)
  selDirection.addEventListener("change", directionPicked)
  selStop.addEventListener("change", stopPicked)
  getRoutes()
}

let selRoute = document.querySelector("#select-route")
let selDirection = document.querySelector("#select-direction")
let selStop = document.querySelector("#select-stop")

loaded()