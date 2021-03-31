world/New()
	rest_server.start()

mob/verb/send_rest_test()
	world << "sending to server"
	var/L = list("status"=0, "players"=0, "admins"=0, "crewManifest"=list("heads"=list(list("name"="Test Dummy", "rank"="Test Dummy"))))
	rest_server.send(L)

	L = list("status"=1, "players"=0, "admins"=0, "crewManifest"=list("heads"=list(list("name"="Test Dummy", "rank"="Test Dummy"))))
	rest_server.send(L)

	L = list("status"=2, "players"=0, "admins"=0, "crewManifest"=list("heads"=list(list("name"="Test Dummy", "rank"="Test Dummy"))))
	rest_server.send(L)

	L = list("status"=3, "players"=0, "admins"=0, "crewManifest"=list("heads"=list(list("name"="Test Dummy", "rank"="Test Dummy"))))
	rest_server.send(L)

	L = list("status"=4, "players"=0, "admins"=0, "crewManifest"=list("heads"=list(list("name"="Test Dummy", "rank"="Test Dummy"))))
	rest_server.send(L)

	L = list("status"=5, "players"=0, "admins"=0, "crewManifest"=list("heads"=list(list("name"="Test Dummy", "rank"="Test Dummy"))))
	rest_server.send(L)
