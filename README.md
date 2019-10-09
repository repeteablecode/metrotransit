# Bus Route

Find the next departure time for a given bus route, desired stop and the routes cardinal direction

## Usage

### Command Format
```
> go run busroute.go "Bus Route" "Bus Stop Name" "Direction"
```

### Sample Output
```
> Departure at 6:15      *No time will show if last stop already departed
```

### Error Message

If no arguments are provided or missing double quotes this error will be presented
```
Arguments Error, Please Format Command as below. Goodbye
   COMMAND: go run casestudy.go "BUS ROUTE" "BUS STOP NAME" "DIRECTION"
   *Note - OS may require escape characters in terminal, example \" and \&
```

### API's used

[MetroTransit](https://svc.metrotransit.org/nextrip/help) API is utilized by this application
