import * as React from "react";
import {type Context, createContext, type ReactNode, useContext, useState} from "react";
import type {ShoppingBasket} from "../domain/shopping-basket/model.tsx";

interface GlobalStateType {
    shoppingBasket: ShoppingBasket
    setShoppingBasket: (shoppingBasket: ShoppingBasket) => void;
}

const GlobalStateContext: Context<GlobalStateType | undefined> = createContext<GlobalStateType | undefined>(undefined);

export const GlobalStateProvider: React.FC<{ children: React.ReactNode }> = ({children}: {
    children: ReactNode
}) => {
    const defaultShoppingBasket : Partial<ShoppingBasket> = {
        items: []
    };
    // @ts-ignore
    const [shoppingBasket, setShoppingBasket] = useState<ShoppingBasket>(defaultShoppingBasket);
    return (
        <GlobalStateContext.Provider value={{setShoppingBasket: setShoppingBasket, shoppingBasket: shoppingBasket}}>
            {children}
        </GlobalStateContext.Provider>
    );
};

export const useGlobalState: () => GlobalStateType = (): GlobalStateType => {
    const context: GlobalStateType | undefined = useContext(GlobalStateContext);
    if (!context) {
        throw new Error("useGlobalState must be used within a GlobalStateProvider");
    }
    return context;
};