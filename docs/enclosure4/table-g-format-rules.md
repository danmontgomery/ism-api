# DoDM 5200.01-V2 Enclosure 4: Table G — Special Format Rules

> **Source**: DoDM 5200.01-V2, February 24, 2012 (Change 4, 7/28/2020), Enclosure 4
>
> **Parent**: [dodm-5200.01-enclosure4-requirements.md](../dodm-5200.01-enclosure4-requirements.md)

---

| ID | Rule | Example | Source |
|----|------|---------|--------|
| **FMT-1** | `//` separates categories | `TOP SECRET//SI//NOFORN` | S1.b.1 |
| **FMT-2** | `/` separates items within same category | `SI/TK` or `NOFORN/PROPIN` | S1.b.1 |
| **FMT-3** | `-` (no spaces) separates control from sub-control | `SI-G`, `HCS-O`, `RD-N` | S1.b.1b |
| **FMT-4** | Space separates sub-compartments | `HCS-O XYZ` | S6.e |
| **FMT-5** | Space separates multiple SIGMA numbers | `RD-SIGMA 1 2` | S8.d.3 |
| **FMT-6** | Comma and space between country codes in REL TO | `REL TO USA, AUS, GBR` | S10.d.4 |
| **FMT-7** | Space between country codes in JOINT | `//JOINT SECRET CAN GBR USA` | S5.e |
| **FMT-8** | Space between country codes in FGI banner | `FGI DEU GBR` | S9.d |
| **FMT-9** | FGI/JOINT documents begin with `//` (no preceding classification) | `//DEU SECRET`, `//JOINT SECRET GBR USA` | S1.c |
| **FMT-10** | U.S. classification is NOT preceded by `//` | `TOP SECRET` (not `//TOP SECRET`) | S3.b |
| **FMT-11** | SAR uses `-` between SAR and nickname/PID | `SAR-BUTTER POPCORN`, `SAR-BP` | S7.c, S7.d |
| **FMT-12** | ACCM uses `-` between ACCM and nickname | `ACCM-FICTITIOUS EFFORT` | S11.a.2 |
| **FMT-13** | 3+ SAPs: banner uses `SAR-MULTIPLE PROGRAMS`; portions must list all PIDs | Banner: `SAR-MULTIPLE PROGRAMS`; Portion: `SAR-TG/SAR-STK/SAR-BP` | S7.e |
