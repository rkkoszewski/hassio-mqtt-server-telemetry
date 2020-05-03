# Server Telemetry Sensor for Home Assistant via MQTT
This is a telemetry service written in Go that publishes information about the server to a MQTT broker.
This service has been designed to be used with a Home Assistant, for which all devices and sensors will show up 
automatically up when autodiscovery is enabled.

## Installation:
0. Install Go and Git on your system (If you haven't yet)
1. Clone this repository
2. In the root of the repository run: *go build -ldflags="-s -w" main.go*
3. Mark the executable as executable (If you're on linux): *chmod +x main*
4. Copy that executable wherever it fits best to you as also thec configuration.yaml file. (Recommended: /opt/server-telemetry/)
5. Edit the configuration.yaml file to your linking.
6. Run the binary and that's it. (By default the binary searches for configuration.yaml, a custom path to the configuration 
file can be set as first argument to the executable)

7. (optional) create service to autostart the script at boot:
    1. sudo nano /etc/systemd/system/system-telemetry-sensor.service
    2. Copy the following content into the service file:
    
    [Unit]
    Description=Home Assistant System Telemetry Sensor service
    After=multi-user.target

    [Service]
    User=[user]
    Type=idle
    ExecStart=/opt/server-telemetry/telemetry-service /opt/server-telemetry/configuration.yaml

    [Install]
    WantedBy=multi-user.target
    
    3. Edit the path to your script path and configuration.yaml. Also make sure you replace [user] with the account from which this script will be run.
    4. sudo systemctl enable system-telemetry-sensor.service 
    5. sudo systemctl start system-telemetry-sensor.service

## To cross compile for AsusWRT
Tested on an Asus RT-AC68U. Before the compilation step set the following environment variables:

set GOOS=linux
set GOARCH=arm
set GOARM=5

## Extra's
You can reduce the size of the binary by commenting out the drivers used in driver/driver.go. This might be improved 
in the future by using compiler flags.

## To Do's
- Automatic removal of disabled sensors is missing. (WORKAROUND: Removing the device from HASS and restarting the service)
- Add support for AsusWRT device temperature measuring (cat /proc/dmu/temperature)
- AsusWRT Device Tracking

## Notes

The project was inspired by: https://github.com/Sennevds/system_sensors