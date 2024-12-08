# Hinoki Planner 🌲

Hinoki Planner is a terminal-based application designed to help you manage your personal life goals effectively.

### Features

- **Timeframes:** Organize your goals using six main timeframes: **day, week, month, quarter, year, and life**. Plan small daily steps and big long-term ambitions simultaneously.
- **Tree Structure:** Link your goals together into a hierarchical tree. Every small daily achievement builds toward your ultimate life goals.
- **Privacy-Focused:** All data is stored locally on your computer, ensuring your plans remain truly personal.
- **Minimal and Powerful:** Hinoki Planner offers a minimalistic interface with smart, hotkey-driven navigation, perfect for users who prioritize robust tools over flashy UI designs.

# Installation
## macOS
```shell
  brew tap kilimondjaro/hinoki-planner
  brew install hinoki-planner 
```

## Linux

# Usage

Hinoki Planner is designed to be highly intuitive and efficient, making navigation quick and easy through a combination of Vim-like keybindings and date shortcuts. With these familiar commands, users can seamlessly manage their goals, navigate through timeframes, and perform tasks without taking their hands off the keyboard.

## Navigation

| Action                                   | Key(s)               | Description                                |
|------------------------------------------|----------------------|--------------------------------------------|
| **Previous period** (e.g., previous day) | `h` or `<-`          | Navigate to the previous period.           |
| **Next period**                          | `l` or `->`          | Navigate to the next period.               |
| **Next goal**                            | `j` or `Arrow Down`  | Move to the next goal.                     |
| **Previous goal**                        | `k` or `Arrow Up`    | Move to the previous goal.                 |
| **Go to date or timeframe**              | `g`                  | Jump to a specific date or timeframe.      |
| **Quit the app**                         | `Q`                  | Exit the application.                      |
| **Select timeframe**                     |                      | Choose the timeframe to navigate to:       |
| Today                                    | `t`                  | Navigate to today.                         |
| Day                                      | `d`                  | Navigate to the current day.               |
| Week                                     | `w`                  | Navigate to the current week.              |
| Month                                    | `m`                  | Navigate to the current month.             |
| Quarter                                  | `q`                  | Navigate to the current quarter.           |
| Year                                     | `y`                  | Navigate to the current year.              |
| Life                                     | `l`                  | Navigate to the entire lifespan timeframe. |

## Managing goals

| Action                            | Key(s)                | Description                                                                                 |
|-----------------------------------|-----------------------|---------------------------------------------------------------------------------------------|
| **Create a goal**                 | `n` then `Enter`      | Create a new goal and submit it by pressing `Enter`.                                        |
| **Mark a goal as done**           | `Spacebar`            | Mark the currently selected goal as done.                                                   |
| **Archive a goal**                | `Backspace`           | Archive the currently selected goal.                                                        |
| **Move a goal to another period** | `D` then specify date | Move the selected goal to another period by pressing uppercase `D` and specifying the date. |
| **Edit a goal**                   | `e`                   | Edit the currently selected goal.                                                           |
| **Open goal details**             | `Enter`               | Open goal details screen for managing subgoals.                                             |

## Date and Timeframe Shortcuts
| Keyword                  | Shorthand     | Description                                                                                   | Timeframe      |
|--------------------------|---------------|-----------------------------------------------------------------------------------------------|----------------|
| `today`                  | `t`           | Current date.                                                                                 | Day            |
| `yesterday`              | `ytd`         | One day before today.                                                                         | Day            |
| `tomorrow`               | `tmrw`        | One day after today.                                                                          | Day            |
| `weekend`                | `wknd`        | Upcoming weekend (Saturday).                                                                  | Day            |
| `monday`                 | `mon`         | Current week's Monday.                                                                        | Day            |
| `tuesday`                | `tue`         | Current week's Tuesday.                                                                       | Day            |
| `wednesday`              | `wed`         | Current week's Wednesday.                                                                     | Day            |
| `thursday`               | `thu`         | Current week's Thursday.                                                                      | Day            |
| `friday`                 | `fri`         | Current week's Friday.                                                                        | Day            |
| `saturday`               | `sat`         | Current week's Saturday.                                                                      | Day            |
| `sunday`                 | `sun`         | Current week's Sunday.                                                                        | Day            |
| `january`                | `jan`         | The month of January.                                                                         | Month          |
| `february`               | `feb`         | The month of February.                                                                        | Month          |
| `march`                  | `mar`         | The month of March.                                                                           | Month          |
| `april`                  | `apr`         | The month of April.                                                                           | Month          |
| `may`                    |               | The month of May.                                                                             | Month          |
| `june`                   | `jun`         | The month of June.                                                                            | Month          |
| `july`                   | `jul`         | The month of July.                                                                            | Month          |
| `august`                 | `aug`         | The month of August.                                                                          | Month          |
| `september`              | `sep`         | The month of September.                                                                       | Month          |
| `october`                | `oct`         | The month of October.                                                                         | Month          |
| `november`               | `nov`         | The month of November.                                                                        | Month          |
| `december`               | `dec`         | The month of December.                                                                        | Month          |
| `week`                   | `w`           | Current week.                                                                                 | Week           |
| `month`                  | `m`           | Current month.                                                                                | Month          |
| `quarter`                | `q`           | Current quarter.                                                                              | Quarter        |
| `q1`                     |               | First quarter of the year.                                                                    | Quarter        |
| `q2`                     |               | Second quarter of the year.                                                                   | Quarter        |
| `q3`                     |               | Third quarter of the year.                                                                    | Quarter        |
| `q4`                     |               | Fourth quarter of the year.                                                                   | Quarter        |
| `year`                   | `y`           | Current year.                                                                                 | Year           |
| `life`                   | `l`           | Represents the entire lifespan timeframe.                                                     | Lifetime       |
| `1 to 31`                |               | A specific day of the current month (e.g., `15` refers to the 15th day of the current month). | Day            |
| `1 to 31 <month>`        |               | A specific day in a specific month of the current year (e.g., `15 March`).                    | Day            |
| `<month> <year>`         |               | A specific month in a specific year (e.g., `March 2024`).                                     | Month          |
| `1 to 31 <month> <year>` |               | A specific day in a specific month in a specific year (e.g., `15 March 2024`).                | Month          |
| `next` + `<key>`         | `n` + `<key>` | Moves to the next occurrence of a timeframe.                                                  | Depends on key |
| `prev` + `<key>`         | `p` + `<key>` | Moves to the previous occurrence of a timeframe.                                              | Depends on key |

# Inspiration
Hinoki Planner draws inspiration from popular todo and planner apps such as **Timestripe, Supernotes, Superlist, Todoist, and Things 3**, among others.

While Hinoki Planner shares some functionality and concepts with these apps, its primary goal is to provide a **minimalist, on-device, terminal-based, hotkey-oriented user experience.** For those who prefer modern UIs and cross-platform support, the apps mentioned above are excellent alternatives.

# License
Hinoki Planner is licensed under the MIT License.
