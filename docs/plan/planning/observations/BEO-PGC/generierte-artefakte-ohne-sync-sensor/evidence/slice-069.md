**Vorgang:** slice-069
**Fund:** Der Reviewer stellte fest, dass es für den neu erzeugten und
committeten Protobuf-Code (`streamv1/*.pb.go`) keinen Sensor gibt, der ihn
gegen die `.proto`-Quelle hält — dieselbe Lücke besteht bereits für
`tools/schema/plan.yaml`/`down.sql`. Er stufte das als INFO ein (kein
Defekt dieses Slice) und empfahl, die Klasse im Register zu führen statt sie
still zu schließen. Der Verifier bestätigte: kein Sync-Gate ist für jetzt
richtig (Kosten pro `make gates`, Repo-Präzedenz), die Bindung bleibt aber
offen und benannt. Erstauftreten.
