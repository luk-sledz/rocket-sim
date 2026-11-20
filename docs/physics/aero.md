# Aerodynamika (2D)

Dokument opisuje model siły aerodynamicznej zaimplementowany w [`internal/rocket`](../../internal/rocket/rocket.go) oraz pochodzenie współczynników użytych w kodzie. Bibliografia: [`docs/refs/README.md`](../refs/README.md).

**Implementacja:** `AttackAngle`, `DragCoeff`, `LiftCoeff`, `DragForce`, `LiftForce`, `AeroForce` w [`rocket.go`](../../internal/rocket/rocket.go). Integrator dodaje `AeroForce` w środku masy wewnątrz `Step` (bez momentu aerodynamicznego).

## Zakres

| Uwzględnione | Poza zakresem |
|--------------|---------------|
| Opór i siła nośna w osiach wiatru (\(C_d(\alpha)\), \(C_l(\alpha)\)) | Wiatr, \(\rho(h)\), tabele Mach / Re |
| \(C_d\) z \(\alpha\) nieznakowanym (`acos`); \(C_l\) z \(\alpha\) znakowanym (`atan2`) | Moment względem środka ciśnienia |
| Stała powierzchnia odniesienia \(S_{\mathrm{ref}}\) | Sterowanie belly-flop / klapy |
| Stała gęstość na poziomie morza | Osie ciała (\(C_N\), \(C_A\)) zamiast osi wiatru |

Logika postawy i ciągu znajduje się w `internal/control` i nie zawiera osobnych członów kompensujących aero.

## Oznaczenia

| Symbol | Definicja |
|--------|-----------|
| `Angle` | Kąt kadłuba [rad]; `0` = dziób w kierunku \(+Y\) |
| \(\hat{u}_{\mathrm{body}}\) | \((\sin\mathrm{Angle},\ \cos\mathrm{Angle})\) |
| \(\mathbf{v}\) / `Vel` | Prędkość inercjalna; przy braku wiatru napływ z \(-\mathbf{v}\) |
| \(\alpha_d\) | Kąt ataku nieznakowany, \(\alpha_d\in[0,\pi]\) — używany w \(C_d\) |
| \(\alpha_l\) | Kąt ataku znakowany, \(\alpha_l\in(-\pi,\pi]\) — używany w \(C_l\) |
| \(S_{\mathrm{ref}}\) | `DragArea` (przekrój spodu / bazy) |

Przy \(|\mathbf{v}|=0\) przyjmuje się \(\alpha_d:=\alpha_l:=0\) oraz \(\mathbf{F}_{\mathrm{aero}}=\mathbf{0}\).

## Równania

### Kąt ataku (opór)

\[
\cos\alpha_d = \hat{u}_{\mathrm{body}}\cdot\hat{u}_{\mathrm{vel}},
\qquad
\alpha_d = \arccos\bigl(\mathrm{clamp}(\cos\alpha_d,-1,1)\bigr).
\]

Odpowiada to `AttackAngle()`. Dla \(C_d\) wystarczy wartość bezwzględna kąta (\(\sin^{2}\alpha_d\)).

### Kąt ataku (nośność)

\[
\sin\alpha_l = \hat{u}_{\mathrm{vel}}\times\hat{u}_{\mathrm{body}}
= u_x b_y - u_y b_x
\quad\text{(składowa 2D)},
\]

\[
\cos\alpha_l = \hat{u}_{\mathrm{body}}\cdot\hat{u}_{\mathrm{vel}},
\qquad
\alpha_l = \mathrm{atan2}(\sin\alpha_l,\ \cos\alpha_l).
\]

Znak \(\alpha_l\) rozróżnia stronę przechyłu \(\pm\mathrm{Angle}\) przy tej samej wartości \(|\alpha|\). Liczony lokalnie w `LiftCoeff()` (nie przez `AttackAngle()`).

### Współczynnik oporu

\[
C_d(\alpha_d) = C_{d0} + (C_{d90}-C_{d0})\,\sin^{2}\alpha_d.
\]

| \(\alpha_d\) | \(\sin^{2}\alpha_d\) | \(C_d\) |
|--------------|----------------------|---------|
| \(0\) lub \(\pi\) (przepływ osiowy) | \(0\) | \(C_{d0}\) |
| \(\pi/2\) (napływ z boku) | \(1\) | \(C_{d90}\) |

### Współczynnik siły nośnej

\[
C_l(\alpha_l)=\tfrac12(C_{d90}-C_{d0})\sin 2\alpha_l.
\]

| \(\alpha_l\) | \(\sin 2\alpha_l\) | \(C_l\) |
|--------------|--------------------|---------|
| \(0\), \(\pm\pi/2\), \(\pm\pi\) | \(0\) | \(0\) |
| \(\pi/4\) | \(+1\) | \(\tfrac12(C_{d90}-C_{d0})\) |
| \(-\pi/4\) lub \(3\pi/4\) | \(-1\) | \(-\tfrac12(C_{d90}-C_{d0})\) |

Przykład (spadek w dół, \(\hat{u}_{\mathrm{vel}}=(0,-1)\), \(\hat{u}_{\perp}=(+1,0)\)):

- \(\mathrm{Angle}=+\pi/4\) → \(\alpha_l=+3\pi/4\), \(C_l<0\) → \(\mathbf{F}_l\) w \(-\!X\) (lewo);
- \(\mathrm{Angle}=-\pi/4\) → \(\alpha_l=-3\pi/4\), \(C_l>0\) → \(\mathbf{F}_l\) w \(+\!X\) (prawo).

### Siła aerodynamiczna

\[
q = \tfrac12\rho\,|\mathbf{v}|^{2},
\]

\[
\mathbf{F}_{d} = -q\,C_d\,S_{\mathrm{ref}}\,\hat{u}_{\mathrm{vel}},
\]

\[
\mathbf{F}_{l} = q\,C_l\,S_{\mathrm{ref}}\,\hat{u}_{\perp},
\qquad
\hat{u}_{\perp}=(-u_{y},\,u_{x}),
\]

\[
\mathbf{F}_{\mathrm{aero}} = \mathbf{F}_{d}+\mathbf{F}_{l}.
\]

`DragForce` i `LiftForce` zwracają osobno \(\mathbf{F}_d\) oraz \(\mathbf{F}_l\); `AeroForce` jest ich sumą.

## Wartości domyślne (rząd wielkości skali Starship)

| Stała / pole | Wartość | Uzasadnienie |
|--------------|---------|--------------|
| `AirDensity` \(\rho\) | \(1.225\,\mathrm{kg\,m^{-3}}\) | Poziom morza, wartość stała |
| `DefaultDragArea` \(S_{\mathrm{ref}}\) | \(63.6\,\mathrm{m^{2}}\) | \(\approx\pi(4.5\,\mathrm{m})^{2}\) (spód) |
| `DefaultDragCd` \(C_{d0}\) | \(1.0\) | Rząd wielkości oporu osiowego / na spód |
| `DefaultDragCd90` \(C_{d90}\) | \(7.08\) | \(C_{d0}\cdot(450/63.6)\); planform \(\approx 50\times 9\,\mathrm{m}\) przy crossflow \(C_{d}\sim 1\) względem planformu, wyrażone względem bazowego \(S_{\mathrm{ref}}\) |

Wartości mają charakter oszacowań inżynierskich, nie dopasowań z tunelu aerodynamicznego ani CFD.

## Pochodzenie postaci \(\sin^{2}\alpha\) oraz \(\sin 2\alpha\)

1. Dane aerodynamiczne smukłych ciał zapisuje się zwykle jako współczynniki względem \(\alpha\) przy **stałym** \(S_{\mathrm{ref}}\). Zmiana „widocznego” kształtu z postawą wchodzi w \(C_d\) oraz \(C_l\), a nie w drugą powierzchnię w kroku całkowania.
2. Przy dużym \(\alpha\) metody inżynierskie dla pocisków i pojazdów nośnych dodają wkład **lepkościowego crossflow** do siły normalnej proporcjonalny do \(C_{dn}\sin^{2}\alpha\) (prędkość poprzeczna \(\sim v\sin\alpha\)). Linia rozwojowa: Allen → Jorgensen (NASA), m.in. [TN D-7228](https://ntrs.nasa.gov/citations/19730012271) oraz [TR R-474](https://ntrs.nasa.gov/citations/19770026166). Pełna lista: [`docs/refs/README.md`](../refs/README.md).
3. Zaimplementowane \(C_d(\alpha_d)\) oraz \(C_l(\alpha_l)\) to **zredukowany model w osiach wiatru**: opór mieszany przez \(\sin^{2}\alpha_d\), nośność przez \(\sin 2\alpha_l\) ze skalą \(\tfrac12(C_{d90}-C_{d0})\). Nie jest to wierna reprodukcja wzorów \(C_N\), \(C_A\), \(C_m\) z TN D-7228.

### Relacja do pełniejszych metod

| Pełniejsza metoda | Ta implementacja |
|-------------------|------------------|
| Osobne \(C_N\) / \(C_A\) (lub \(C_d\)+\(C_l\)) | \(C_d(\alpha_d)\) oraz \(C_l(\alpha_l)\) w osiach wiatru |
| \(C_{dn}(M_n,Re_n)\), współczynnik smukłości \(\eta\) | Stałe \(C_{d0}\), \(C_{d90}\) |
| Moment aerodynamiczny względem COP | Siła wyłącznie w COM |
| Atmosfera i wiatr | Stałe \(\rho\); inercjalne `Vel` |

## Mapowanie na kod

| Wielkość | Symbol / API |
|----------|--------------|
| \(\alpha_d\) | `AttackAngle()` |
| \(C_d(\alpha_d)\) | `DragCoeff()` |
| \(\alpha_l\), \(C_l(\alpha_l)\) | `LiftCoeff()` |
| \(\mathbf{F}_d\) | `DragForce()` |
| \(\mathbf{F}_l\) | `LiftForce()` |
| \(\mathbf{F}_{\mathrm{aero}}\) | `AeroForce()` = `DragForce` + `LiftForce` |
| Całkowanie | `Step` → `force.Add(r.AeroForce())` |
| Testy regresji | `TestAeroForceZeroVel`, `TestAeroForceOpposesVel`, `TestDragCoeffVsAttack`, `TestLiftCoeffVsAttack` |

## Prace przyszłe

- Prędkość względna \(\mathbf{v}_{\mathrm{rel}}=\mathbf{v}-\mathbf{W}\).
- Moment aerodynamiczny (COP ≠ COM).
- \(\rho(h)\) oraz tabelaryczne współczynniki.
- Model w osiach ciała (\(C_n\sim\sin^{2}\alpha\)) zamiast (lub obok) osi wiatru.

Przy każdej zmianie równań aerodynamicznych w kodzie należy zaktualizować ten dokument w tej samej zmianie.
