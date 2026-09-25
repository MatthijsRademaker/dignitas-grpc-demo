# Third-party notices

## Golden duck artwork

`src/art/duck.ts`, `src/components/GoldenDuck.vue` and `public/favicon.svg` use the shapes of the duck
emoji from **Noto Emoji**:

- Project: <https://github.com/googlefonts/noto-emoji>
- File: `2D/svg/emoji_u1f986.svg` (commit `9411589`)
- Licence: Apache License, Version 2.0

We changed it: the original fills are replaced by gold gradients, and we added an outline, specular highlights
and an animated shimmer. The shapes themselves are unchanged. Notice from `2D/svg/LICENSE`:

```
Copyright 2013 Google, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Fonts

Bundled into the build from the `@fontsource` npm packages. Each one is licensed under the SIL Open Font
License 1.1, and the full licence text ships in the package's `LICENSE` file.

| Font | Copyright | Package |
| --- | --- | --- |
| Fraunces | 2020 The Fraunces Project Authors | `@fontsource-variable/fraunces` |
| Nunito Sans | 2016 The Nunito Sans Project Authors | `@fontsource-variable/nunito-sans` |
| Fira Code | 2014-2020 The Fira Code Project Authors | `@fontsource/fira-code` |
