# Subnet

Le subnet ti permettono di segmentare la rete VPC in sezioni più piccole e isolate. Ogni subnet può avere il proprio blocco CIDR e può essere usata per organizzare risorse, controllare il routing e applicare politiche di sicurezza all'interno di un VPC.

## Comandi

### Elenca Subnet

Elenca tutte le subnet in un VPC.

```bash
acloud network subnet list <vpc-id>
```

**Argomenti:**
- `vpc-id` - L'ID del VPC

**Esempio:**
```bash
acloud network subnet list 689307f4745108d3c6343b5a
```

**Output:**
```
NAME         ID                        CIDR           STATUS
subnet-1     1234567890abcdef          10.0.1.0/24    Active
subnet-2     0987654321fedcba          10.0.2.0/24    Active
```

### Ottieni Dettagli Subnet

Ottieni dettagli su una subnet specifica.

```bash
acloud network subnet get <vpc-id> <subnet-id>
```

**Argomenti:**
- `vpc-id` - L'ID del VPC
- `subnet-id` - L'ID della subnet

**Esempio:**
```bash
acloud network subnet get 689307f4745108d3c6343b5a 1234567890abcdef
```

**Output:**
```
Subnet Details:
===============
ID:            1234567890abcdef
Name:          subnet-1
CIDR:          10.0.1.0/24
Status:        Active
```

### Crea Subnet

Crea una nuova subnet in un VPC.

```bash
acloud network subnet create <vpc-id> --cidr <cidr> --name <name>
```

**Flag Richiesti:**
- `--cidr string` - Blocco CIDR per la subnet
- `--name string` - Nome per la subnet

**Esempio:**
```bash
acloud network subnet create 689307f4745108d3c6343b5a --cidr 10.0.1.0/24 --name subnet-1
```

**Output:**
```
Subnet created successfully!
ID:      1234567890abcdef
Name:    subnet-1
CIDR:    10.0.1.0/24
```

### Aggiorna Subnet

Aggiorna il nome o il CIDR di una subnet esistente.

```bash
acloud network subnet update <vpc-id> <subnet-id> [flags]
```

**Argomenti:**
- `vpc-id` - L'ID del VPC
- `subnet-id` - L'ID della subnet

**Flag:**
- `--name string` - Nuovo nome per la subnet
- `--cidr string` - Nuovo blocco CIDR

**Esempio:**
```bash
acloud network subnet update 689307f4745108d3c6343b5a 1234567890abcdef --name new-subnet-name
```

**Output:**
```
Subnet updated successfully!
ID:      1234567890abcdef
Name:    new-subnet-name
```

### Elimina Subnet

Elimina una subnet da un VPC.

```bash
acloud network subnet delete <vpc-id> <subnet-id>
```

**Argomenti:**
- `vpc-id` - L'ID del VPC
- `subnet-id` - L'ID della subnet

**Esempio:**
```bash
acloud network subnet delete 689307f4745108d3c6343b5a 1234567890abcdef
```

**Output:**
```
Subnet 1234567890abcdef deleted successfully!
```

## Auto-completamento Shell

I comandi subnet supportano auto-completamento per ID VPC e ID subnet.

## Best Practices

- Usa nomi descrittivi per le subnet in base al loro scopo.
- Evita blocchi CIDR sovrapposti tra subnet.

## Risoluzione dei Problemi

- Assicurati che il VPC sia **Active** prima di creare subnet.
- Controlla conflitti CIDR quando aggiungi o aggiorni subnet.

## Comandi Correlati

- [VPC](vpc.md) - Gestisci VPC
- [Security Group](securitygroup.md) - Gestisci security group
