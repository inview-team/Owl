from opcua import Server
import random
import datetime

server = Server()

url = "opc.tcp://127.0.0.1:30329"
server.set_endpoint(url)

name = "OPC_UA_SERVER"
addspace = server.register_namespace(name)

list_of_nodes = []

for i in range(0, 8):
    list_of_nodes.append(server.get_objects_node())

list_of_params = []

for i in range(0, 8):
    list_of_params.append(list_of_nodes[i].add_object(addspace, "Parameters"))

Press = list_of_params[0].add_variable(addspace, "Pressure", 0)
Humid = list_of_params[1].add_variable(addspace, "Humidity", 0)
roomTemp = list_of_params[2].add_variable(addspace, "Room temperature", 0)
workingAreaTemp = list_of_params[3].add_variable(addspace, "Temperature of the working area", 0)
pH = list_of_params[4].add_variable(addspace, "Level of pH", 0)
weight = list_of_params[5].add_variable(addspace, "Weight", 0)
fluidFl = list_of_params[6].add_variable(addspace, "Fluid flow", 0)
co2 = list_of_params[7].add_variable(addspace, "Level of CO2", 0)

Press.set_writable() #ns=2;i=9
Humid.set_writable() #ns=2;i=10
roomTemp.set_writable() #ns=2;i=11
workingAreaTemp.set_writable() #ns=2;i=12
pH.set_writable()   #ns=2;i=13
weight.set_writable() #ns=2;i=14
fluidFl.set_writable() #ns=2;i=15
co2.set_writable() #ns=2;i=16

server.start()
print("Server started at {}".format(url))

while True:
    for i in range(0, 100):
        Pressure = random.uniform(900.0, 1200.0)
        Humidity = random.uniform(0.0, 100.0)
        roomTemperature = random.uniform(-40.0, 100)
        workingAreaTemperature = random.uniform(100.0, 1000.0)
        levelOfpH = random.uniform(0.0, 14.0)
        Weight = random.uniform(0.0, 1000.0)
        fluidFlow = random.uniform(0.0, 100.0)
        CO2 = random.uniform(0.0, 100.0)
        TIME = datetime.datetime.now()


        print(datetime.datetime.now(), Pressure, Humidity, roomTemperature, workingAreaTemperature, levelOfpH, Weight, fluidFlow, CO2, TIME)

        Press.set_value(Pressure)
        Humid.set_value(Humidity)
        roomTemp.set_value(roomTemp)
        workingAreaTemp.set_value(workingAreaTemperature)
        pH.set_value(levelOfpH)
        weight.set_value(Weight)
        fluidFl.set_value(fluidFlow)
        co2.set_value(CO2)

