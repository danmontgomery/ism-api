# DoDM 5200.01-V2 Enclosure 4: Comprehensive Example Markings for Test Validation

> **Source**: DoDM 5200.01-V2, February 24, 2012 (Change 4, 7/28/2020), Enclosure 4
>
> **Parent**: [dodm-5200.01-enclosure4-requirements.md](../dodm-5200.01-enclosure4-requirements.md)

---

## Valid Banner Lines

| # | Banner Line | Key Features |
|---|-------------|--------------|
| 1 | `TOP SECRET` | Simple US classification |
| 2 | `SECRET` | Simple US classification |
| 3 | `CONFIDENTIAL` | Simple US classification |
| 4 | `TOP SECRET//SI//NOFORN` | SCI + dissemination |
| 5 | `SECRET//SI/TK//RELIDO` | Multiple SCI + dissemination |
| 6 | `TOP SECRET//SI-GAMMA//ORCON/NOFORN` | SCI with compartment + multiple dissem |
| 7 | `TOP SECRET//HCS//NOFORN` | HCS requires NOFORN |
| 8 | `SECRET//TK-GEOCAP//NOFORN` | TK-GEOCAP requires NOFORN |
| 9 | `TOP SECRET//SAR-BUTTERED POPCORN` | SAP with nickname |
| 10 | `TOP SECRET//SAR-MULTIPLE PROGRAMS` | 3+ SAPs |
| 11 | `SECRET//RESTRICTED DATA` | RD |
| 12 | `SECRET//RESTRICTED DATA-N` | CNWDI |
| 13 | `SECRET//RD-SIGMA 1 2` | Multiple SIGMAs |
| 14 | `SECRET//FRD-SIGMA 14` | FRD with SIGMA |
| 15 | `TOP SECRET//FGI DEU GBR` | FGI in US doc |
| 16 | `SECRET//FGI NATO` | NATO info in US doc |
| 17 | `TOP SECRET//REL TO USA, EGY, ISR` | REL TO |
| 18 | `SECRET//NOFORN` | Document with NOFORN/REL TO mix (NOFORN wins) |
| 19 | `SECRET//DISPLAY ONLY AFG` | DISPLAY ONLY |
| 20 | `SECRET//REL TO USA, GBR/DISPLAY ONLY AFG` | REL TO + DISPLAY ONLY (all portions match) |
| 21 | `SECRET//ORCON//NOFORN` | Separator between categories |
| 22 | `TOP SECRET//ORCON//NOFORN` | ORCON with NOFORN example |
| 23 | `SECRET//ACCM-FICTITIOUS EFFORT/TEA LEAF` | ACCM with multiple programs |
| 24 | `SECRET//IMCON` | IMCON (SECRET only) |
| 25 | `SECRET//IMCON/RELIDO` | IMCON with RELIDO |
| 26 | `TOP SECRET//IMCON/NOFORN` | IMCON + NOFORN in mixed doc |
| 27 | `//JOINT SECRET CAN GBR USA` | JOINT marking |
| 28 | `//JOINT SECRET GBR USA//REL TO USA, AUS, CAN, GBR, NZL` | JOINT + REL TO |
| 29 | `//DEU SECRET` | Non-US FGI document |
| 30 | `//COSMIC TOP SECRET` | NATO TS |
| 31 | `//COSMIC TOP SECRET BOHEMIA` | NATO TS SIGINT |
| 32 | `//NATO SECRET` | NATO SECRET |
| 33 | `CONFIDENTIAL//NOFORN/PROPIN` | NOFORN + PROPIN |
| 34 | `TOP SECRET//NOFORN/FISA` | NOFORN + FISA |
| 35 | `SECRET//EXDIS` | DoS EXDIS |
| 36 | `SECRET//NODIS` | DoS NODIS |
| 37 | `CONFIDENTIAL//SI//REL TO USA, AUS, FRA` | SCI with REL TO |
| 38 | `TOP SECRET//SI-GAMMA//SAR-BP//NOFORN` | SCI + SAP + NOFORN |
| 39 | `SECRET//FORMERLY RESTRICTED DATA` | FRD |
| 40 | `TOP SECRET//SAR-TIN BAKER//WAIVED` | SAP WAIVED |
| 41 | `TOP SECRET//TK//SAR-BP` | SCI + SAP in correct order |
| 42 | `SECRET//HCS-O XYZ//NOFORN` | HCS with compartment and sub-compartment |

## Invalid Banner Lines (Should Fail Validation)

| # | Invalid Banner Line | Violation | Source |
|---|---------------------|-----------|--------|
| 1 | `SECRET//NOFORN/REL TO USA, GBR` | NOFORN and REL TO in same banner | MX-1 |
| 2 | `SECRET//NOFORN/RELIDO` | NOFORN and RELIDO in same banner | MX-2 |
| 3 | `SECRET//EXDIS/NODIS` | EXDIS and NODIS together | MX-4 |
| 4 | `SECRET//DISPLAY ONLY AFG/NOFORN` | DISPLAY ONLY with NOFORN | MX-6 |
| 5 | `SECRET//DISPLAY ONLY AFG/RELIDO` | DISPLAY ONLY with RELIDO | MX-7 |
| 6 | `CONFIDENTIAL//IMCON` | IMCON requires SECRET (not CONFIDENTIAL) | CL-11 |
| 7 | `TOP SECRET//IMCON` | IMCON standalone requires SECRET | CL-11 |
| 8 | `SECRET//REL TO USA` | REL TO USA alone is not approved | CR-3 |
| 9 | `TOP SECRET//HCS` | HCS without NOFORN | CR-1 |
| 10 | `SECRET//TK-GEOCAP` | TK-GEOCAP without NOFORN | CR-2 |
| 11 | `//JOINT SECRET GBR USA` then `TOP SECRET` | JOINT and US classification mixed | MX-9 |
| 12 | `SECRET//SAR-BP//REL TO USA, GBR` then `SECRET//SAR-BP//NOFORN` | Mixed NOFORN/REL TO portions should produce NOFORN banner | BP-1 |
| 13 | `TOP SECRET//RESTRICTED DATA` then `UNCLASSIFIED//RD` | RD with UNCLASSIFIED | CL-1 |
| 14 | `//COSMIC TOP SECRET BOHEMIA` at `//NATO SECRET` | BOHEMIA only with COSMIC TOP SECRET | CR-9 |
| 15 | `SECRET//FOUO` | FOUO should not appear in classified banner | CR-11 |
| 16 | `SECRET//REL TO GBR, USA, AUS` | USA must be first in REL TO list | CC-1 |
| 17 | `//JOINT SECRET USA CAN GBR` | JOINT countries must be alphabetical (USA not first) | CC-2 |
| 18 | `SECRET/SI//NOFORN` | Single slash between categories (should be `//`) | FMT-1 |
| 19 | `SECRET//SI GAMMA//NOFORN` | Space instead of hyphen for compartment | FMT-3 |
| 20 | `SECRET//SAR-BP//SAR-MDP//SAR-TG` | 3+ SAPs should use MULTIPLE PROGRAMS in banner | E4-S7.e |
| 21 | `UNCLASSIFIED//RELIDO` | RELIDO requires TS, S, or C | CL-10 |
| 22 | `UNCLASSIFIED//ORCON` | ORCON requires TS, S, or C | CL-5 |
| 23 | `SECRET//ACCM-FE` | ACCM must use full nickname, no abbreviations | E4-S11.a.5 |
| 24 | `//JOINT RESTRICTED GBR USA` | RESTRICTED not allowed when US is co-owner | CL-15 |
