# DoDM 5200.01-V2 Enclosure 4: JOINT Classification Markings

> **Source**: DoDM 5200.01-V2, February 24, 2012 (Change 4, 7/28/2020), Enclosure 4, Section 5
>
> **Parent**: [dodm-5200.01-enclosure4-requirements.md](../dodm-5200.01-enclosure4-requirements.md)

---

## 5.1 Format

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S5.c** | All JOINT markings (banner and portion) begin with `//JOINT` | S5.c |
| **E4-S5.d** | Required banner format: `//JOINT [classification] [country codes]` | S5.d |
| **E4-S5.d.restr** | If US is NOT a co-owner, classification may be RESTRICTED. Where US IS a co-owner, RESTRICTED may NOT be used. | S5.d |
| **E4-S5.e** | Country codes (including USA) listed in alphabetical order, separated by spaces | S5.e |
| **E4-S5.e.note** | USA placement in JOINT is alphabetical (unique to JOINT; differs from REL TO where USA is first) | S5.e |

## 5.2 Portion Marking Rules

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S5.f** | Country codes are NOT included in portion markings when all portions match the banner country codes | S5.f |
| **E4-S5.f.extract** | If a JOINT portion is extracted into a non-JOINT U.S. document, country codes MUST be listed alphabetically in portion markings: `(//JOINT [classification] [country codes])` | S5.f |

## 5.3 Derivative Document Rules

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S5.g** | JOINT portions must be segregated from U.S. classified information in derivative documents | S5.g |
| **E4-S5.g.banner** | Banner of derivative U.S. document uses highest classification of all portions as a U.S. classification | S5.g |
| **E4-S5.g.nojoint** | JOINT marking is NOT carried to the banner line; used only in applicable portions | S5.g |
| **E4-S5.g.fgi** | FGI markings shall be added to the banner line with all non-U.S. country codes from JOINT portions | S5.g |

## 5.4 JOINT with REL TO

| ID | Requirement | Example Banner | Example Portion | Source |
|----|-------------|----------------|-----------------|--------|
| **E4-S5.fig31** | JOINT information may be combined with REL TO | `//JOINT SECRET GBR USA//REL TO USA, AUS, CAN, GBR, NZL` | `(//JOINT S//REL)` | Figure 31 |
| **E4-S5.fig31.rel** | `(REL)` may be used in portion if REL TO country list matches banner line | | `(//JOINT S//REL)` | Figure 31 |

## 5.5 Classification Authority

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S5.h** | Classification authority block is used ONLY when the United States is one of the co-owners | S5.h |
