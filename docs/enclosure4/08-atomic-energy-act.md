# DoDM 5200.01-V2 Enclosure 4: Atomic Energy Act (AEA) Information Markings

> **Source**: DoDM 5200.01-V2, February 24, 2012 (Change 4, 7/28/2020), Enclosure 4, Section 8
>
> **Parent**: [dodm-5200.01-enclosure4-requirements.md](../dodm-5200.01-enclosure4-requirements.md)

---

## 8.1 Restricted Data (RD)

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S8.a.4** | RD may be used ONLY with TOP SECRET, SECRET, or CONFIDENTIAL | S8.a.(4) |
| **E4-S8.a.5** | Banner: `[classification]//RESTRICTED DATA` (or `RD`). Portion: `([classification]//RD)` | S8.a.(5) |
| **E4-S8.a.5a** | If any portion contains RD, RD must appear in the banner line | S8.a.(5) |
| **E4-S8.a.7** | RD is not subject to automatic declassification; "Declassify On:" shall state "Not applicable" or be deleted | S8.a.(7) |
| **E4-S8.a.10** | Do not commingle RD and NSI in the same portion | S8.a.(10)(a) |

## 8.2 Formerly Restricted Data (FRD)

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S8.b.2** | FRD may be used ONLY with TOP SECRET, SECRET, or CONFIDENTIAL | S8.b.(2) |
| **E4-S8.b.3** | Banner: `[classification]//FORMERLY RESTRICTED DATA` (or `FRD`). Portion: `([classification]//FRD)` | S8.b.(3) |
| **E4-S8.b.3a** | If any portion contains FRD, FRD must appear in the banner line | S8.b.(3) |
| **E4-S8.b.5** | FRD is not subject to automatic declassification; "Declassify On:" shall state "Not applicable" or be deleted | S8.b.(5) |
| **E4-S8.b.8.a** | Do not commingle FRD and NSI in the same portion | S8.b.(8)(a) |

## 8.3 Critical Nuclear Weapons Design Information (CNWDI)

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S8.c.1** | CNWDI is the designation for TOP SECRET RD or SECRET RD weapons data | S8.c.(1) |
| **E4-S8.c.1.class** | CNWDI applies only to TOP SECRET or SECRET (as subset of RD) | S8.c.(1) |
| **E4-S8.c.3** | Append `-N` to banner and portion markings (e.g., `SECRET//RESTRICTED DATA-N`, portion `(S//RD-N)`) | S8.c.(3) |

## 8.4 SIGMA

| ID | Requirement | Source |
|----|-------------|--------|
| **E4-S8.d.2** | SIGMA may be used ONLY with TOP SECRET, SECRET, or CONFIDENTIAL | S8.d.(2) |
| **E4-S8.d.3** | Banner format: `[classification]//[RD or FRD]-SIGMA [#]`. Sigma number 1-99. | S8.d.(3) |
| **E4-S8.d.3.multi** | Multiple SIGMAs listed in numerical order, separated by a space (e.g., `RD-SIGMA 1 2`) | S8.d.(3) |
| **E4-S8.d.3.portion** | Portion: use `SG` with sigma number (e.g., `(S//RD-SG 1)`, `(S//FRD-SG 14)`) | S8.d.(3) |

## 8.5 AEA Examples

| Banner Line | Portion Marking | Notes |
|-------------|----------------|-------|
| `SECRET//RESTRICTED DATA` | `(S//RD)` | Basic RD |
| `SECRET//FORMERLY RESTRICTED DATA` | `(S//FRD)` | Basic FRD |
| `SECRET//RESTRICTED DATA-N` | `(S//RD-N)` | CNWDI |
| `SECRET//RD-SIGMA 1 2` | `(S//RD-SG 1)`, `(S//RD-SG 2)` | Multiple SIGMA |
| `SECRET//FRD-SIGMA 14` | `(S//FRD-SG 14)` | FRD with SIGMA |
