# DoDM 5200.01-V2 Enclosure 4: SAP Control Markings

> **Source**: DoDM 5200.01-V2, February 24, 2012 (Change 4, 7/28/2020), Enclosure 4, Section 7
>
> **Parent**: [dodm-5200.01-enclosure4-requirements.md](../dodm-5200.01-enclosure4-requirements.md)

---

## 7.1 Banner Format

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S7.c** | Banner: classification, caveat "SAR" or "SPECIAL ACCESS REQUIRED", and nickname/code word with dissemination control if assigned | S7.c |
| **E4-S7.c.hyphen** | Hyphen (`-`) without spaces separates "SAR" and the program nickname/codeword | S7.c |
| **E4-S7.c.slash** | Forward slash (`/`) separates multiple nicknames/codewords | S7.c |
| **E4-S7.c.noPID** | Assigned program identifiers (PIDs) shall NOT be used in the banner line | S7.c |

## 7.2 Portion Format

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S7.d** | Portion: classification, SAR, and the program's assigned PID (e.g., `TS//SAR-BP`) | S7.d |
| **E4-S7.d.hyphen** | Hyphen without spaces separates SAR and PID | S7.d |
| **E4-S7.d.multi** | Multiple PIDs listed alphabetically, separated by a slash | S7.d |

## 7.3 Multiple Programs

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S7.e** | When 3+ SAPs in a single document, use `SAR-MULTIPLE PROGRAMS` in banner (e.g., `SECRET//SAR-MULTIPLE PROGRAMS`) | S7.e |
| **E4-S7.e.portions** | When there are multiple programs, the PID for each SAP MUST be cited in the portion marking regardless of total number of PIDs | S7.e |
| **E4-S7.e.alpha** | Multiple PIDs in portion markings must be listed alphabetically, separated by `/` | S7.e |

## 7.4 WAIVED

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S7.f** | SAPs exempted from Congressional reporting: marked "WAIVED" in banner, portion marks, and on media (e.g., `TOP SECRET//SAR-BP//WAIVED`) | S7.f |
| **E4-S7.f.dissem** | "WAIVED" is used as a dissemination control marking | S7.f |

## 7.5 HVSACO

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S7.g** | HVSACO (Handle via Special Access Channels Only) is a handling instruction, NOT a classification level or dissemination control | S7.g |
| **E4-S7.g.apply** | Applied to non-SAP material (classified or unclassified) within a SAP environment (e.g., `SECRET//HVSACO`, `UNCLASSIFIED//HVSACO`) | S7.g |

## 7.6 SAP Examples (Figure 34)

| Banner Line | Portion Marking | Notes |
|-------------|----------------|-------|
| `TOP SECRET//SPECIAL ACCESS REQUIRED-BUTTERED POPCORN` | `(TS//SAR-BP)` | Full name in banner, PID in portion |
| `TOP SECRET//SAR-SWAGGER` | `(TS//SAR-SGR)` | Acronym OK in banner |
| `TOP SECRET//TALENT KEYHOLE//SAR-BP` | `(TS//TK//SAR-BP)` | SCI + SAP |
| `TOP SECRET//SAR-BLUE FROG/SAR-MUDDY PATH` | `(TS//SAR-BFG/SAR-MDP)` | Two programs |
| `TOP SECRET//SAR-MULTIPLE PROGRAMS` | `(TS//SAR-TG/SAR-STK/SAR-BP)` | 3+ programs: MULTIPLE PROGRAMS in banner |
| `TOP SECRET//SAR-TIN BAKER//WAIVED` | `(TS//SAR-TB//WAIVED)` | WAIVED dissemination |
| `TOP SECRET//SAR-DAGGER//WAIVED` | `(TS//SAR-DGR//WAIVED)` | WAIVED dissemination |
