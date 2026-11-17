# Referencje

Źródła koncepcji crossflow / smukłych ciał przy dużym kącie ataku, wykorzystanych w [`docs/physics/aero.md`](../physics/aero.md). Preferowane są **strony cytowań NTRS**; pliki PDF należy pobierać z poziomu strony raportu (bezpośrednie adresy `/api/.../downloads/` bywają zawodne).

## Raporty techniczne NASA

1. **Jorgensen, L. H.** — NASA TN D-7228 (1973)  
   *A method for estimating static aerodynamic characteristics for slender bodies of circular and noncircular cross section alone and with lifting surfaces at angles of attack from 0° to 90°*  
   - [NTRS](https://ntrs.nasa.gov/citations/19730012271)  
   - Wykorzystanie: wkład lepkościowego crossflow \(\propto C_{dn}\sin^{2}\alpha\) do siły normalnej. Bez implementacji pełnej metody \(C_N\), \(C_A\), \(C_m\).

2. **Jorgensen, L. H.** — NASA TR R-474 (1977)  
   *Prediction of static aerodynamic characteristics for slender bodies alone and with lifting surfaces to very high angles of attack*  
   - [NTRS](https://ntrs.nasa.gov/citations/19770026166)

3. **Procedura inżynierska w linii Allena** — NASA TN D-6996  
   *Prediction of static aerodynamic characteristics for space-shuttle-like and other bodies at angles of attack from 0° to 180°*  
   - [NTRS](https://ntrs.nasa.gov/citations/19730006261)  
   - Kontekst: człon potencjalny oraz empiryczny lepkościowy crossflow (linia Allena leżąca u podstaw metod Jorgensena).

## Uzupełnienie dydaktyczne

4. MIT OpenCourseWare 2.154 / 13.49 — wykład 10 (zwarte zestawienie formuł typu Jorgensen)  
   - [Zasób lec10](https://ocw.mit.edu/courses/2-154-maneuvering-and-control-of-surface-and-underwater-vehicles-13-49-fall-2004/resources/lec10/)  
   - [Strona kursu](https://ocw.mit.edu/courses/2-154-maneuvering-and-control-of-surface-and-underwater-vehicles-13-49-fall-2004/)  
   - Materiał pomocniczy; do cytowania kanonicznego służą raporty NASA powyżej.

## Licencja / atrybucja

Raporty NASA NTRS są utworami Rządu Stanów Zjednoczonych dostępnymi do użytku publicznego. W cytowaniu należy podawać numer raportu (np. TN D-7228). Repozytorium implementuje uproszczony model MVP pod soft catch i nie deklaruje zgodności jeden-do-jednego z przywołanymi metodami.
