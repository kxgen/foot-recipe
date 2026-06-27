# Tailwind CSS Cheat Sheet: The 50 Most Common Classes

This reference guide covers the 50 most frequently used Tailwind CSS utility classes, categorized by function. It provides their exact vanilla CSS equivalents and clear descriptions to help bridge the gap between utility-first classes and traditional CSS.

---

## 📐 Layout & Box Model

### 1. Flexbox & Grid
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `flex` | `display: flex;` | Establishes a flex container, enabling a flex context for its immediate children. |
| `grid` | `display: grid;` | Establishes a grid container, enabling a grid context for its immediate children. |
| `hidden` | `display: none;` | Completely removes the element from the document layout (hides it). |
| `block` | `display: block;` | Displays an element as a block-level element, starting on a new line. |
| `inline-block`| `display: inline-block;` | Displays an element as an inline-level block container (flows with text but honors width/height). |

### 2. Spacing (Margin & Padding)
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `m-4` | `margin: 1rem; /* 16px */` | Applies uniform margin on all four sides of the element. |
| `mx-auto` | `margin-left: auto; margin-right: auto;` | Centers a block element horizontally within its parent container. |
| `mt-2` | `margin-top: 0.5rem; /* 8px */` | Applies margin only to the top side of the element. |
| `p-4` | `padding: 1rem; /* 16px */` | Applies uniform inner padding on all four sides of the element. |
| `px-6` | `padding-left: 1.5rem; padding-right: 1.5rem;` | Applies padding simultaneously to the left and right sides. |
| `py-3` | `padding-top: 0.75rem; padding-bottom: 0.75rem;` | Applies padding simultaneously to the top and bottom sides. |

### 3. Sizing
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `w-full` | `width: 100%;` | Forces the element's width to span the entirety of its parent container. |
| `w-screen` | `width: 100vw;` | Sets the element's width to span 100% of the viewport width. |
| `h-full` | `height: 100%;` | Sets the element's height to 100% of its parent's height. |
| `h-screen` | `height: 100vh;` | Sets the element's height to span 100% of the viewport height. |
| `max-w-md` | `max-width: 28rem; /* 448px */` | Imposes a hard limit on how wide an element can expand. |

---

## 🧩 Alignment & Positioning

### 4. Flexbox Alignment
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `items-center` | `align-items: center;` | Aligns flex items along the cross axis (vertically if flex-row). |
| `justify-center`| `justify-content: center;` | Centers flex items along the main axis (horizontally if flex-row). |
| `justify-between`| `justify-content: space-between;` | Distributes flex items evenly; first item is at start, last is at end. |
| `flex-col` | `flex-direction: column;` | Stack flex items vertically instead of horizontally. |
| `gap-4` | `gap: 1rem; /* 16px */` | Controls the space between rows and columns in flexbox or grid. |

### 5. Standard Positioning
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `absolute` | `position: absolute;` | Positions an element relative to its nearest positioned ancestor. |
| `relative` | `position: relative;` | Positions an element relative to its normal document flow position. |
| `fixed` | `position: fixed;` | Positions an element relative to the viewport window; stays in place during scroll. |
| `top-0` | `top: 0px;` | Anchors a positioned element to the very top edge of its boundary box. |
| `z-50` | `z-index: 50;` | Controls the 3D stack order layer of an element (brings it to front). |

---

## 🎨 Visual Styling (Typography, Colors, Borders)

### 6. Typography
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `text-sm` | `font-size: 0.875rem; line-height: 1.25rem;`| Sets font size to small and regulates line spacing. |
| `text-lg` | `font-size: 1.125rem; line-height: 1.75rem;`| Sets font size to large and regulates line spacing. |
| `text-2xl` | `font-size: 1.5rem; line-height: 2rem;` | Sets font size to extra-large (ideal for headings). |
| `font-bold` | `font-weight: 700;` | Formats text with a heavy, bold line weight. |
| `font-medium` | `font-weight: 500;` | Formats text with a medium line weight (slightly bolder than regular). |
| `text-center` | `text-align: center;` | Aligns running text into the direct horizontal center of its box. |
| `text-white` | `color: rgb(255 255 255);` | Changes text color to solid pure white. |
| `text-gray-700` | `color: rgb(55 65 81);` | Changes text color to a balanced dark grey color. |

### 7. Backgrounds
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `bg-white` | `background-color: rgb(255 255 255);` | Fills the background of the element with solid pure white. |
| `bg-blue-500` | `background-color: rgb(59 130 246);` | Fills the background with a vibrant corporate blue hue. |
| `bg-gray-100` | `background-color: rgb(243 244 246);` | Fills the background with a soft, clean light grey accent. |
| `bg-transparent`| `background-color: transparent;` | Makes the background completely clear and see-through. |

### 8. Borders & Curves
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `border` | `border-width: 1px;` | Renders a standard single-pixel border line around the element. |
| `border-gray-300`| `border-color: rgb(209 213 219);` | Tint the border color with a neutral subtle light grey. |
| `rounded` | `border-radius: 0.25rem; /* 4px */` | Slightly rounds the sharp corners of an element's border box. |
| `rounded-md` | `border-radius: 0.375rem; /* 6px */` | Applies a moderate, modernized corner curve to an element. |
| `rounded-full` | `border-radius: 9999px;` | Fully rounds corners to produce pill shapes or perfect circles. |

---

## ✨ Effects, Interactions & States

### 9. Depth & Opacity
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `shadow` | `box-shadow: 0 1px 3px 0 rgba(0,0,0,0.1), ...;`| Applies a soft drop shadow beneath the element for depth. |
| `shadow-lg` | `box-shadow: 0 10px 15px -3px rgba(0,0,0,0.1), ...;`| Applies a prominent, deep drop shadow for overlay elements. |
| `opacity-50` | `opacity: 0.5;` | Makes the element 50% translucent / semi-transparent. |

### 10. Transitions & States
| Tailwind Class | Equivalent Vanilla CSS | Description |
| :--- | :--- | :--- |
| `transition` | `transition-property: color, background-color...;`| Prepares elements to change styles smoothly instead of snapping instantly.|
| `duration-300` | `transition-duration: 300ms;` | Tells transitions to take exactly 300 milliseconds to complete. |
| `hover:bg-blue-600`| `&:hover { background-color: rgb(37 99 235); }`| Modifies the background color only when the mouse cursor is hovering over it.|
| `focus:outline-none`| `&:focus { outline: 2px solid transparent; }`| Strips out the default browser ring focus outline when active. |
| `cursor-pointer` | `cursor: pointer;` | Changes the mouse mouse-cursor into a hand pointer shape (indicates clickability). |

---

*Tip: You can pair any of these classes with responsive modifiers (like `md:flex` or `lg:text-2xl`) to build fully responsive designs with zero custom media queries!*