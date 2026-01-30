## Step 9: (Opzionale) Crea un file user-data per cloud-init

Se vuoi personalizzare il server all'avvio (ad esempio, creare utenti o installare pacchetti), crea un file `cloud-init.yaml` prima del provisioning. Esempio:

```yaml
#cloud-config
users:
	- name: demo
		sudo: ALL=(ALL) NOPASSWD:ALL
		groups: users, admin
		home: /home/demo
		shell: /bin/bash
		lock_passwd: false
		passwd: <hashed-password>
```

Sostituisci `<hashed-password>` con una password hash valida (vedi la documentazione cloud-init per dettagli). Puoi aggiungere altre opzioni cloud-init secondo necessità.

## Step 10: Attendi che il Cloud Server sia Attivo

Dopo aver eseguito il comando di creazione, puoi controllare lo stato del server con:

```bash
acloud compute cloudserver list | grep "my-server"
```

Quando lo stato è `Active`, il server è pronto all'uso:

```
my-server                 697c62605b79733376b3386a       ITBG-Bergamo    CSO4A8          Active
```

---

## Step 11: Connettersi al Cloud Server

Per connetterti al cloud server, usa il comando `acloud compute cloudserver connect` specificando l'utente in base all'immagine di avvio. Per Ubuntu, usa l'utente `ubuntu`:

```bash
acloud compute cloudserver connect 697c62605b79733376b3386a --user ubuntu
```

Esempio output:
```
Connect by running: ssh ubuntu@85.235.152.94
```

> **Nota:**
> - Il flag `--user` è obbligatorio e deve corrispondere all'utente predefinito dell'immagine (es. `ubuntu` per Ubuntu, `centos` per CentOS, ecc.).
> - Il comando connect mostrerà il comando SSH da utilizzare.

---

## Step 12: Spegnere o Accendere il Cloud Server

Puoi spegnere o accendere il cloud server in qualsiasi momento con i seguenti comandi:

Per spegnere:

```bash
acloud compute cloudserver power-off 697c62605b79733376b3386a
```

Esempio output:
```
Cloud server powered off successfully!
Server: my-server
Status: Updating
```

Per accendere:

```bash
acloud compute cloudserver power-on 697c62605b79733376b3386a
```

Esempio output:
```
Cloud server powered on successfully!
Server: my-server
Status: Updating
```

---

## Step 13: Eliminare il Cloud Server

Per eliminare il cloud server quando non serve più, usa il comando:

```bash
acloud compute cloudserver delete 697c62605b79733376b3386a
```

Ti verrà chiesta conferma:

```
Are you sure you want to delete cloud server '697c62605b79733376b3386a'? (yes/no): yes
Cloud server '697c62605b79733376b3386a' deleted successfully.
```
# Esempio di Provisioning Cloud Server

Questo esempio mostra come eseguire il provisioning di un nuovo cloud server utilizzando la CLI ArubaCloud, con tutti i flag di rete richiesti e le istruzioni aggiornate. Segui la versione inglese per i dettagli dei passaggi e aggiorna i valori secondo il tuo ambiente.

Per la guida dettagliata, consulta la documentazione in inglese nella stessa directory.
