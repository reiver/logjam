import { ThirdwebAuth } from "@thirdweb-dev/auth/next";
import { PrivateKeyWallet } from "@thirdweb-dev/auth/evm";

export const { ThirdwebAuthHandler, getUser } = ThirdwebAuth({
  wallet: new PrivateKeyWallet(import.meta.env.THIRDWEB_AUTH_PRIVATE_KEY || ""),
  domain: "localhost:3000",
});
