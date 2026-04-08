# DoDM 5200.01-V2 Enclosure 4: SCI Control System Markings

> **Source**: DoDM 5200.01-V2, February 24, 2012 (Change 4, 7/28/2020), Enclosure 4, Section 6
>
> **Parent**: [dodm-5200.01-enclosure4-requirements.md](../dodm-5200.01-enclosure4-requirements.md)

---

## 6.1 Published SCI Control Systems

| ID | System | Abbreviation | Source |
|----|--------|-------------|--------|
| **E4-S6.b.1** | HUMINT Control System | HCS | S6.b.(1) |
| **E4-S6.b.2** | Special Intelligence | SI | S6.b.(2) |
| **E4-S6.b.3** | TALENT KEYHOLE | TK | S6.b.(3) |

## 6.2 SCI Formatting Rules

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S6.c** | Multiple SCI entries are listed alphabetically, separated by a single forward slash (`/`) | S6.c |
| **E4-S6.d** | A hyphen without spaces separates an SCI control system name and its compartment(s) (e.g., `SI-GAMMA`, `SI-G-XXX`) | S6.d |
| **E4-S6.d.multi** | Multiple compartments are listed alpha-numerically | S6.d |
| **E4-S6.e** | Sub-compartments are separated from their compartment and each other by a space (` `); listed alpha-numerically | S6.e |

## 6.3 SCI Co-requirements

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S6.f** | When HCS is used, NOFORN must also be used | S6.f |
| **E4-S6.f.tk** | When TK-GEOCAP is used, NOFORN must also be used | S6.f |

## 6.4 SCI Processing Constraints

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S6.g** | SCI, regardless of classification level, must be processed only on SCI-accredited systems (e.g., JWICS), NOT on SIPRNET | S6.g |

## 6.5 SCI Banner/Portion Examples (Figure 33)

| Banner Line | Portion Marking | Notes |
|-------------|----------------|-------|
| `TOP SECRET//HCS//NOFORN` | `(TS//HCS//NF)` | HCS requires NOFORN |
| `SECRET//SI/TK//RELIDO` | `(S//SI/TK//RELIDO)` | Multiple SCI, alphabetical |
| `TOP SECRET//SI-GAMMA//ORCON/NOFORN` | `(TS//SI-G//OC/NF)` | SCI with compartment, multiple dissem |
| `CONFIDENTIAL//SI//REL TO USA, AUS, FRA` | `(C//SI//REL TO USA, FRA)` | SCI with REL TO |
| `TOP SECRET//SI-XXX//REL TO USA, AUS` | `(TS//SI-XXX//REL)` | Portion uses (REL) when same as banner |
| `SECRET//HCS-O XYZ//NOFORN` | `(S//HCS-O XYZ//NF)` | Compartment with sub-compartment |
| `SECRET//TK-GEOCAP//NOFORN` | `(S//TK-G//NF)` | TK-GEOCAP requires NOFORN |
