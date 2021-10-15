import { Fider } from "@fider/services"
import { FiderImpl } from "../fider"

export const fiderMock = {
  notAuthenticated: (): FiderImpl => {
    return Fider.initialize({
      settings: {
        environment: "development",
        oauth: [],
        ldap: [],
      },
      tenant: {},
      user: undefined,
    })
  },
  authenticated: (): FiderImpl => {
    return Fider.initialize({
      settings: {
        environment: "development",
        oauth: [],
        ldap: [],
      },
      tenant: {},
      user: {
        name: "Jon Snow",
      },
    })
  },
}
