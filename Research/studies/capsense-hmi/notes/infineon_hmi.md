# Infineon CapSense та HMI поза automotive

> **WP type:** Analysis Report — supporting research note<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/research_notes/Застосування CapSense поза авто/infineon_hmi.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/research_notes/%D0%97%D0%B0%D1%81%D1%82%D0%BE%D1%81%D1%83%D0%B2%D0%B0%D0%BD%D0%BD%D1%8F%20CapSense%20%D0%BF%D0%BE%D0%B7%D0%B0%20%D0%B0%D0%B2%D1%82%D0%BE/infineon_hmi.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


Дослідження за офіційними матеріалами Infineon/Cypress, перевіреними 11 листопада 2026 року. Сторінки Infineon здебільшого не вказують дату редакції; де версія/дата відсутня, це зазначено. Під «CapSense» тут мається на увазі CAPSENSE™ у MCU PSoC, а не окремий каталог capacitive sensor IC, зорієнтований переважно на автомобільні задачі.

## Які можливості CapSense придатні для неавтомобільних виробів?

### Takeaway
CapSense у сучасному PSoC — це програмована підсистема вимірювання ємності й HMI-віджетів на електродах плати: кнопки, матриці, слайдери, touchpad/трекпад, proximity та жести. PSoC 4000T додає п’яте покоління CapSense, низьке енергоспоживання, рідинну толерантність і Multi-Sense (зокрема індуктивне та безконтактне вимірювання рівня рідини). Це не означає, що MCU сам по собі є повним контролером графічного сенсорного екрана.

### Cited Findings
- Офіційний репозиторій CAPSENSE Middleware перелічує методи CSD (self-capacitance), CSX (mutual-capacitance) та ISX (inductive), а серед widget-типів — buttons, matrix buttons, sliders, touchpads і proximity sensors; також згадано one-/two-finger gestures і liquid-level measurement. README має copyright 2019–2026, без номера окремої редакції. [Infineon CAPSENSE Middleware README](https://github.com/Infineon/capsense/blob/master/README.md)
- Сторінка PSoC 4000T вказує п’яте покоління CapSense, 16 self-capacitance і 64 mutual-capacitance входи, до 8 індуктивних входів, hover touch і deep-sleep 6 µA. Це характеристики сімейства, а не твердження, що всі функції одночасно доступні в кожній конфігурації чи на будь-якому наборі виводів. [PSoC 4000T](https://www.infineon.com/products/microcontroller/32-bit-psoc-arm-cortex/psoc-4-mcu/4000/4000t)
- Infineon описує PSoC 4000T Multi-Sense як поєднання ємнісного, індуктивного й вимірювання рівня рідини; названі інтерфейси включають touch, slider, force-touch, proximity, finger tracking, metal touch і hover touch. [PSoC 4000T](https://www.infineon.com/products/microcontroller/32-bit-psoc-arm-cortex/psoc-4-mcu/4000/4000t)
- Загальна сторінка PSoC 4 характеризує CAPSENSE як вбудований периферійний блок для кнопок, слайдерів, proximity та вимірювання рівня рідини; також окремо від CapSense пояснює Multi-Sense з індуктивним sensing. [PSoC 4 MCU](https://www.infineon.com/products/microcontroller/32-bit-psoc-arm-cortex/psoc-4-mcu)
- Для реалізації конкретного дизайну Infineon посилається на application note AN85951, *PSoC 4 and PSoC 6 MCU CAPSENSE™ Design Guide*, та на ModusToolbox CAPSENSE Configurator/Tuner guides. Сторінка PSoC 4000T не показує номер редакції чи дату AN85951; Configurator і Tuner — інструменти конфігурації/налаштування, не готові UI-контролери. [AN85951](https://www.infineon.com/content/dam/infineon/row/public/documents/faq/an85951-psoc-r-4-and-psoc-6-mcu-capsense-r-design-guide.pdf) · [Configurator User Guide](https://www.infineon.com/ModusToolboxCapSenseConfig) · [Tuner Tool Guide](https://www.infineon.com/ModusToolboxCapSenseTuner)

### Inferences
- PSoC із CapSense природно підходить, коли потрібні сенсорні електроди, події/жести, MCU-логіка та інші вбудовані функції в одному виробі. Потрібно звірити sensor count, доступні GPIO, режими CSD/CSX, scan time, електроди й firmware конкретного part number.
- Touchpad widget або mutual-capacitance входи не дорівнюють гарантованому XY-контролеру для дисплея. Відображення графіки, стек дисплея, повноцінний touchscreen controller і потрібний протокол мають бути розглянуті як окремі функції/компоненти, якщо цього вимагає виріб.

### Gaps
- Розглянуті сторінки не задають універсальні розміри електродів, товщину панелі, точність XY або максимальну кількість одночасних точок для довільного сенсорного екрана. Ці параметри потребують дизайну з AN85951 і перевірки конкретного MCU та stack-up.

## Які неавтомобільні приклади, набори й застосування називає Infineon?

### Takeaway
HMI — центральний і явно маркетований напрям для PSoC CapSense поза automotive. Первинні приклади Infineon охоплюють wearables/hearables, побутову техніку та кухонні прилади, димери й touchpad-керування навушниками; доступні набори демонструють кнопки, slider, proximity і touchpad. Промислова автоматизація теж присутня в описі набору CY8CKIT-145-40XX, але в переглянутих матеріалах не знайдено окремого медичного reference design чи заяви про медичну кваліфікацію.

### Cited Findings
- Сторінка *Capacitive Touch Control* прямо позиціонує CapSense поза авто, називає AR glasses і home appliances, а також gloves, liquid tolerance, touch buttons, air-gap touch, capacitive sliders, trackpads і customizable gestures. [Infineon Capacitive Touch Control](https://www.infineon.com/applications/solutions/human-machine-interface/touch-control)
- Та сама сторінка наводить окремі прикладні вертикалі: wearables (XY touchpad, gestures, wear detection), dimming controllers (finger tracking, gesture, thick overlay), headphones (напрямні жести, play/pause, volume/finger tracking) і home appliances/RAC/induction cooktops; для appliance прикладу заявлено роботу через air gap «up to 10 mm». Це демонстраційні застосування/можливості сторінки, а не гарантована відстань для будь-якого матеріалу чи електрода. [Infineon Capacitive Touch Control](https://www.infineon.com/applications/solutions/human-machine-interface/touch-control)
- Активна сторінка набору CY8CKIT-040T (PSoC 4000T Evaluation Kit) перелічує кнопку, proximity sensor і touchpad та позначає застосування refrigerator, fitness bands і capacitive touch control. Infineon описує демонстрацію: плата може бути занурена до позначки immersible line або обприскана краплями без хибних touch-подій; touchpad разом із proximity sensor імітує wearable screen із bezel. Сторінка на момент перевірки має статус Active and preferred; сторінка не датує свою редакцію. [CY8CKIT-040T product page](https://www.infineon.com/evaluation-board/CY8CKIT-040T) · [CY8CKIT-040T User Guide](https://www.infineon.com/assets/row/public/documents/30/44/infineon-cy8ckit-040t-user-guide-usermanual-en.pdf)
- CY8CKIT-041S-MAX (PSoC 4100S Max) демонструє CapSense touchpad, slider, buttons і proximity sensor; сторінка зазначає LCD drive capability і 5th Gen CapSense. Сам факт LCD drive не є доказом наявності проєкційно-ємнісного графічного touchscreen controller. [CY8CKIT-041S-MAX](https://www.infineon.com/evaluation-board/CY8CKIT-041S-MAX)
- CY8CKIT-145-40XX — активний набір PSoC 4000S із трьома CapSense-кнопками й 5-сегментним slider. На сторінці вказані target applications: home entertainment, industrial automation, wearables; кваліфікація виробу позначена як Industrial, Consumer. [CY8CKIT-145-40XX](https://www.infineon.com/evaluation-board/CY8CKIT-145-40XX)
- Infineon позначає сторінку CY8CKIT-040T як актуальну оцінювальну платформу, але це демонстраційний kit, а не закончений appliance HMI reference design із сертифікованою панеллю. Його user guide доступний зі сторінки набору. [CY8CKIT-040T](https://www.infineon.com/evaluation-board/CY8CKIT-040T)

### Inferences
- Для швидкого PoC вибір залежить від сценарію: CY8CKIT-145-40XX — базові кнопки/slider; CY8CKIT-041S-MAX — ширше MCU-середовище з touchpad/slider/proximity; CY8CKIT-040T — нова liquid-tolerance/proximity демонстрація. Спершу звірити доступність, сумісність інструментів і потрібні електроди в актуальному kit guide.
- Для industrial HMI продуктова можливість реалістична, але набір чи загальна сторінка застосувань не підтверджують сертифікацію готового промислового пристрою.
- Медичний інтерфейс можна розглядати як окрему інженерну можливість, але у переглянутих офіційних джерелах не знайдено явного medical CapSense reference design, медичного призначення чи твердження про медичну сертифікацію. Не варто наводити medical як підтверджений Infineon use case без додаткового продуктового джерела.

### Gaps
- Виявлені приклади на сторінці touch-control — застосувальні описи, а не всі повні схеми, BOM, layout, firmware чи звіти про випробування. CY8CKIT-040T є доказом демонстраційних характеристик конкретного набору, але не серійної панелі refrigerator.
- Не знайдено окремого актуального офіційного медичного CapSense design/reference design у переглянутому матеріалі.

## Де межа між CapSense touchpad і touchscreen?

### Takeaway
У джерелах PSoC CapSense має touchpad/trackpad widget і може виконувати finger tracking чи жести на дискретних електродах. Це придатне для панелі керування, та не є автоматично повним projected-capacitive (PCAP) touchscreen-контролером для дисплея.

### Cited Findings
- Middleware визначає touchpad як один із CAPSENSE widgets поряд із кнопками, матрицями, slider і proximity. [CAPSENSE Middleware README](https://github.com/Infineon/capsense/blob/master/README.md)
- Infineon називає у touch-control матеріалі trackpads і XY touchpad інтерфейси для окремих виробів (зокрема wearables і headphones), а не заявляє, що PSoC 4000T є дисплейним PCAP controller. [Capacitive Touch Control](https://www.infineon.com/applications/solutions/human-machine-interface/touch-control)
- Kit CY8CKIT-040T окремо називає свою touchpad із proximity sensor емуляцією «wearable screen with a bezel» — це формулювання стосується демонстрації сенсорного інтерфейсу, не вбудованого екрана. [CY8CKIT-040T](https://www.infineon.com/evaluation-board/CY8CKIT-040T)
- PSoC 4100S Max kit перераховує LCD drive capability разом із CapSense-входами, однак product page не визначає це як PCAP display controller. [PSoC 4100S Max kit](https://www.infineon.com/evaluation-board/CY8CKIT-041S-MAX)

### Inferences
- Для механічно окремої панелі з кнопками/slider/невеликим трекпадом CapSense природний вибір. Для touch поверх LCD/OLED із багатоточковим XY на всій площі потрібен окремий аналіз сенсорного стеку, контролера, дисплея/інтерфейсу та підтримуваних розмірів; не слід вибирати PSoC лише тому, що він має touchpad widget або багато mutual-capacitance входів.
- «Touchpad» описує функцію для координат/жестів з електродів, «touchscreen» — сенсорний дисплейний модуль/система; терміни не взаємозамінні без конкретної схеми та даташиту.

### Gaps
- У розглянутих матеріалах не було достатньо даних, щоб оцінити підтримку PSoC 4000T або іншого конкретного PSoC як дисплейного PCAP контролера. Заявляти таку придатність без спеціального контролерного даташиту/референс дизайну не можна.

## Що обмежує water/glove/overlay/proximity claims і як це стосується дока?

### Takeaway
Infineon заявляє glove touch, рідинну толерантність, thick overlay й air-gap operation, але це властивості налаштованого sensing-дизайну, а не blanket guarantee для довільного складання. У розглянутому dock-дослідженні CapSense важливіший як можливий допоміжний proximity/occupancy-вхід; сенсорна панель — необов’язковий HMI для налаштування, сервісу чи ручного керування, не центральний механізм посадки/виявлення дрона.

### Cited Findings
- Infineon окремо заявляє glove touch, liquid tolerance, operation through air gaps, operation through thick overlays і різні HMI widgets. Сторінка не визначає товщину чи матеріал рукавички/панелі, за яких усі ці властивості гарантовані. [Capacitive Touch Control](https://www.infineon.com/applications/solutions/human-machine-interface/touch-control)
- CY8CKIT-040T має конкретний демонстраційний тест: занурення до позначки на платі або краплі без false touches. Це не є IP-рейтингом, доказом герметичності кінцевого виробу чи тестом для довільної води/солі/миючих засобів/конденсату. [CY8CKIT-040T](https://www.infineon.com/evaluation-board/CY8CKIT-040T)
- Сторінка PSoC 4000T описує liquid tolerance у мокрих/high-moisture умовах і називає highly conductive salty water; також заявляє proximity detection up to 30 cm. Джерело не задає універсальні габарити електрода/цілі чи stack-up для таких чисел. Не переносити 30 см на саморобний корпус без перевірки. [PSoC 4000T](https://www.infineon.com/products/microcontroller/32-bit-psoc-arm-cortex/psoc-4-mcu/4000/4000t)
- Технічним керівництвом Infineon для конкретної конструкції є AN85951 CAPSENSE Design Guide; поточна сторінка PSoC 4000T також посилається на AN234231 з методами зниження споживання. Сторінка не надає продуктового IP рейтингу чи сертифікації готового HMI-виробу. [AN85951](https://www.infineon.com/content/dam/infineon/row/public/documents/faq/an85951-psoc-r-4-and-psoc-6-mcu-capsense-r-design-guide.pdf) · [AN234231](https://www.infineon.com/row/public/documents/30/42/infineon-an234231-psoc4-capsense-ultralowpower-capacitive-sensing-techniques-applicationnotes-en.pdf)
- Зовнішня категорія Infineon «Capacitive sensors» описана переважно в контексті automotive: hands-on detection, AEC-Q100, functional-safety/ASIL. Не слід плутати її з вбудованою PSoC CAPSENSE підсистемою чи використовувати automotive qualification заяву як підтвердження кваліфікації неавтомобільного HMI. [Infineon Capacitive Sensors](https://www.infineon.com/products/sensor/capacitive-sensors)
- У вихідних нотатках цього дослідження CapSense було оцінено як спосіб отримувати ємнісний proximity/occupancy сигнал біля landing-box; джерело цієї конкретної архітектури — попередній робочий контекст, не Infineon reference design. [Попередні нотатки CapSense для dock](../../sensor-options/notes/capsense.md)

### Inferences
- Для заяви «працює в рукавичках» треба задати тип/товщину рукавички, положення руки, матеріал і товщину overlay, electrode geometry, grounded/metal parts і провести вимірювання на зібраному продукті. Офіційна сторінка підтверджує capability, не кожне поєднання умов.
- Для води треба розділяти (1) вологу/краплі, (2) тимчасове занурення сенсорної плати, (3) ступінь захисту корпусу та (4) безпечність/надійність кінцевої системи. Демонстрація одного kit не засвідчує інші три.
- Товщина й діелектричні властивості накладки, повітряний зазор, заземлення, екранування, паразитна ємність і шум визначають чутливість та хибні спрацьовування. Формулювання «до 10 мм» чи «до 30 см» треба трактувати як заяву про конкретні демонстраційні/налаштовані умови, не загальну гарантію.
- Для дока CapSense може виявляти зміну ємнісного зв’язку в налаштованій зоні; сам по собі він не визначає клас об’єкта, не вимірює повітряний потік і не доводить посадку дрона. Touch HMI може бути доданий для power/status/mode/reset/service, але це відокремлена функція від датчиків dock-позиції. Тож у цій роботі HMI — peripheral scope, якщо немає окремої вимоги до локального інтерфейсу.
- Сертифікацію медичного чи промислового готового виробу, вимоги safety, EMC, IP/environmental, ризики операторського інтерфейсу та software validation треба підтверджувати для повного продукту. Переглянуті CapSense-кити й маркетингові сторінки цього не засвідчують.

### Gaps
- Немає специфікації майбутнього dock overlay, електродів, матеріалів корпусу, рівня води/забруднень, рукавичок, металевих деталей чи потрібної дистанції; отже, неможливо стверджувати, що HMI чи proximity забезпечать потрібну точність/надійність.
- У розглянутому офіційному матеріалі не знайдено CapSense reference design для docking box, drone landing, медичного сертифікованого HMI чи повного графічного touchscreen.
- Перш ніж приймати вимоги «glove/water/overlay», слід вибрати конкретний PSoC і перевірити актуальні datasheet, AN85951 та kit/user guide, потім прототипувати реальний stack-up і провести EMC, false-touch і environmental tests.
