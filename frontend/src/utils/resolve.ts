type Resolved<Data> = [undefined, Data] | [Error, undefined]

export const resolve = <Data>(promise: Promise<Data>): Promise<Resolved<Data>> =>
  promise
    .then((data): [undefined, Data] => [undefined, data])
    .catch((error): [Error, undefined] => [error, undefined])
