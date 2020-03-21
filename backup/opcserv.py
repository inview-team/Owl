import sys
from opcua import ua, Server
from trumania.core.random_generators import NumpyRandomGenerator

server = Server()

server.set_endpoint("opc.tcp://0.0.0.0:4840/")
server.set_server_name("Server")
uri = "metrics"
idx = server.register_namespace(uri)
objects = server.get_objects_node()

metricTypes = ["Pressure", "Humidity", "RoomTemp", "WorkTemp", "FluidFlow", "Mass", "PH", "CO2"]
objects = []
gens = []
for mt in metricTypes:
    objects.append(objects.add_object(idx,mt))
    gen.append(NumpyRandomGenerator(method="random_sample(100)", loc=3, scale=5, seed=next(example_circus.seeder)))

Discret_1 = Object_1.add_variable(idx,'Discret_1',[0,0,0,0,0,0,0,0])  
Discret_2 = Object_2.add_variable(idx,'Discret_2',[0,0,0,0,0,0,0,0])  
Analog_3  = Object_3.add_variable(idx,'Analog_3',[10,20,30,40,50])  

server.start() 
